package handler

import (
	"context"
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/application"
	"github.com/iamnalinor/vk-mattermost-polls/internal/models"
	"github.com/iamnalinor/vk-mattermost-polls/internal/repo"
	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func getPollHandler(ctx context.Context, post *model.Post) (bool, error) {
	args := strings.Split(post.Message, " ")
	if len(args) != 2 {
		return true, RespondToPost(ctx, post, "Usage: `!poll <id>`")
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return true, RespondToPost(ctx, post, "ID must be integer")
	}

	app := ctx.Value("app").(*application.Application)
	poll, err := app.Repos.Poll.GetByID(id)
	if repo.IsNotFound(err) {
		return true, RespondToPost(ctx, post, "Poll does not exist. Possibly, it was deleted by its creator")
	}
	if err != nil {
		return true, fmt.Errorf("get poll: %w", err)
	}

	app.Logger.Debug("user viewed poll",
		zap.String("channelId", post.ChannelId),
		zap.String("userId", post.UserId),
		zap.Int("pollId", id))

	votes, err := app.Repos.PollVote.GetByPollID(id)
	if err != nil {
		return true, fmt.Errorf("get poll votes: %w", err)
	}
	votesCount := make(map[int]int)
	for _, vote := range votes {
		votesCount[vote.OptionIndex]++
	}

	text := fmt.Sprintf("Poll: %s\n\nOptions:\n", poll.Question)
	for i, option := range poll.Options {
		text += fmt.Sprintf("%d. %s (%d votes)\n", i+1, option, votesCount[i])
	}
	if poll.Status == models.PollActive {
		text += fmt.Sprintf("\nVote via `!vote %d <option number>`\n", id)
		if post.UserId == poll.CreatorID {
			text += fmt.Sprintf("Close poll: `!closepoll %d`\n", id)
		}
	} else {
		text += "\nPoll is closed\n"
	}

	if post.UserId == poll.CreatorID {
		text += fmt.Sprintf("Delete poll: `!deletepoll %d`\n", id)
	}

	return true, RespondToPost(ctx, post, text)
}
