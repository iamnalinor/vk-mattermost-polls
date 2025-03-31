package repo

import (
	"errors"
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/models"
	"github.com/tarantool/go-tarantool/v2"
)

const pollSpace = "poll"

type PollRepo struct {
	conn *tarantool.Connection
}

func NewPollRepo(conn *tarantool.Connection) (*PollRepo, error) {
	req := tarantool.NewEvalRequest(`
if box.space[...] == nil then
    box.schema.space.create(..., {
        if_not_exists = true,
        format = {
            { name = "id", type = "unsigned" },
            { name = "creator_id", type = "string" },
            { name = "question", type = "string" },
            { name = "options", type = "array" },
            { name = "status", type = "string" },
        }
    })
    box.space[...]:create_index("primary", {
        type = "TREE",
        parts = { { field = 1, type = "unsigned" } },
        sequence = true,
        if_not_exists = true,
    })
end
    `).Args([]any{pollSpace, pollSpace, pollSpace})
	if _, err := conn.Do(req).Get(); err != nil {
		return nil, fmt.Errorf("create poll space: %w", err)
	}

	return &PollRepo{conn: conn}, nil
}

func (r *PollRepo) GetByID(id int) (models.Poll, error) {
	req := tarantool.NewSelectRequest(pollSpace).Index("primary").Key([]any{id})
	data, err := r.conn.Do(req).Get()

	if err != nil {
		var clientErr tarantool.ClientError
		if errors.As(err, &clientErr) && clientErr.Code == notFoundErrorCode {
			return models.Poll{}, ErrNotFound
		}
		return models.Poll{}, err
	}
	if len(data) == 0 {
		return models.Poll{}, ErrNotFound
	}

	row := data[0].([]any)
	return models.Poll{
		ID:        asInt(row[0]),
		CreatorID: row[1].(string),
		Question:  row[2].(string),
		Options:   asStringSlice(row[3]),
		Status:    models.PollStatus(row[4].(string)),
	}, nil
}

func (r *PollRepo) Create(poll models.Poll) (int, error) {
	req := tarantool.NewInsertRequest(pollSpace).Tuple([]any{
		nil,
		poll.CreatorID,
		poll.Question,
		poll.Options,
		poll.Status,
	})
	data, err := r.conn.Do(req).Get()
	if err != nil {
		return 0, err
	}

	id := data[0].([]any)[0]
	return asInt(id), nil
}

func (r *PollRepo) Update(poll models.Poll) error {
	req := tarantool.NewReplaceRequest(pollSpace).Tuple([]any{
		poll.ID,
		poll.CreatorID,
		poll.Question,
		poll.Options,
		poll.Status,
	})
	_, err := r.conn.Do(req).Get()
	return err
}

func (r *PollRepo) Delete(id int) error {
	req := tarantool.NewDeleteRequest(pollSpace).Index("primary").Key([]any{id})
	_, err := r.conn.Do(req).Get()
	return err
}
