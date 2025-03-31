package models

type PollStatus string

const (
	PollActive PollStatus = "active"
	PollClosed PollStatus = "closed"
)

type Poll struct {
	ID        int
	CreatorID string
	Question  string
	Options   []string
	Status    PollStatus
}
