package repo

import (
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/models"
	"github.com/tarantool/go-tarantool/v2"
)

type Poll interface {
	GetByID(id int) (models.Poll, error)
	Create(poll models.Poll) (int, error)
	Update(poll models.Poll) error
	Delete(id int) error
}

type PollVote interface {
	GetByPollID(id int) ([]models.PollVote, error)
	Upsert(pollVote models.PollVote) error
}

type Repositories struct {
	Poll     Poll
	PollVote PollVote
}

func NewRepositories(conn *tarantool.Connection) (*Repositories, error) {
	pollRepo, err := NewPollRepo(conn)
	if err != nil {
		return nil, fmt.Errorf("init poll repo: %w", err)
	}
	pollVoteRepo, err := NewPollVoteRepo(conn)
	if err != nil {
		return nil, fmt.Errorf("init poll vote repo: %w", err)
	}
	return &Repositories{
		Poll:     pollRepo,
		PollVote: pollVoteRepo,
	}, nil
}
