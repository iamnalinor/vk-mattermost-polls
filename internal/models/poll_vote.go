package models

type PollVote struct {
	PollID      int
	UserID      string
	OptionIndex int
}
