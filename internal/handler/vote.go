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

func voteHandler(ctx context.Context, post *model.Post) (bool, error) {
	args := strings.Split(post.Message, " ")
	if len(args) != 3 {
		return true, RespondToPost(ctx, post, "Usage: `!vote <poll id> <option number>`")
	}
	pollId, err := strconv.Atoi(args[1])
	if err != nil {
		return true, RespondToPost(ctx, post, "ID must be integer")
	}
	option, err := strconv.Atoi(args[2])
	if err != nil {
		return true, RespondToPost(ctx, post, "Option must be integer")
	}

	app := ctx.Value("app").(*application.Application)
	poll, err := app.Repos.Poll.GetByID(pollId)
	if repo.IsNotFound(err) {
		return true, RespondToPost(ctx, post, "Poll does not exist. Possibly, it was deleted by its creator")
	}
	if err != nil {
		return true, fmt.Errorf("get poll: %w", err)
	}

	if poll.Status == models.PollClosed {
		return true, RespondToPost(ctx, post, "Sorry, poll is closed")
	}

	if option < 1 || option > len(poll.Options) {
		return true, RespondToPost(ctx, post, "Option does not exist")
	}

	if err := app.Repos.PollVote.Upsert(models.PollVote{
		PollID:      pollId,
		UserID:      post.UserId,
		OptionIndex: option - 1,
	}); err != nil {
		return true, fmt.Errorf("upsert poll vote: %w", err)
	}

	app.Logger.Info("user voted", zap.Int("pollId", pollId), zap.String("userId", post.UserId), zap.Int("option", option))

	return true, RespondToPost(ctx, post, fmt.Sprintf("You voted for option %s", poll.Options[option-1]))
}
