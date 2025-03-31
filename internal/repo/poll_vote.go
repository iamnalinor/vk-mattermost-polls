package repo

import (
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/models"
	"github.com/tarantool/go-tarantool/v2"
)

const pollVoteSpace = "poll_vote"

type PollVoteRepo struct {
	conn *tarantool.Connection
}

func NewPollVoteRepo(conn *tarantool.Connection) (*PollVoteRepo, error) {
	req := tarantool.NewEvalRequest(`
if box.space[...] == nil then
    box.schema.space.create(..., {
        if_not_exists = true,
        format = {
            { name = "poll_id", type = "unsigned" },
            { name = "user_id", type = "string" },
            { name = "option_index", type = "unsigned" },
        }
    })
    box.space[...]:create_index('primary', {
        parts = {'poll_id', 'user_id'},
        if_not_exists = true,
    })
    box.space[...]:create_index('poll_id', {
        parts = {'poll_id'},
        unique = false,
        if_not_exists = true,
	})
end
    `).Args([]any{pollVoteSpace, pollVoteSpace, pollVoteSpace, pollVoteSpace})
	if _, err := conn.Do(req).Get(); err != nil {
		return nil, fmt.Errorf("create poll space: %w", err)
	}
	return &PollVoteRepo{conn: conn}, nil
}

func (r *PollVoteRepo) GetByPollID(id int) ([]models.PollVote, error) {
	req := tarantool.NewSelectRequest(pollVoteSpace).Index("poll_id").Key([]any{id})
	data, err := r.conn.Do(req).Get()
	if err != nil {
		return nil, err
	}
	var votes []models.PollVote
	for _, row := range data {
		row2 := row.([]any)
		votes = append(votes, models.PollVote{
			PollID:      asInt(row2[0]),
			UserID:      row2[1].(string),
			OptionIndex: asInt(row2[2]),
		})
	}
	return votes, nil
}

func (r *PollVoteRepo) Upsert(pollVote models.PollVote) error {
	req := tarantool.NewUpsertRequest(pollVoteSpace).
		Tuple([]any{pollVote.PollID, pollVote.UserID, pollVote.OptionIndex}).
		Operations(tarantool.NewOperations().Assign(2, pollVote.OptionIndex))
	_, err := r.conn.Do(req).Get()
	return err
}
