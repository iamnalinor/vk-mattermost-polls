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

func closePollHandler(ctx context.Context, post *model.Post) (bool, error) {
	args := strings.Split(post.Message, " ")
	if len(args) != 2 {
		return true, RespondToPost(ctx, post, "Usage: `!closepoll <id>`")
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return true, RespondToPost(ctx, post, "ID must be integer")
	}

	app := ctx.Value("app").(*application.Application)
	poll, err := app.Repos.Poll.GetByID(id)
	if repo.IsNotFound(err) {
		return true, RespondToPost(ctx, post, "Poll does not exist")
	}
	if err != nil {
		return true, fmt.Errorf("get poll: %w", err)
	}

	if poll.CreatorID != post.UserId {
		return true, RespondToPost(ctx, post, "You are not the creator of this poll")
	}

	if poll.Status == models.PollClosed {
		return true, RespondToPost(ctx, post, "Poll is already closed")
	}

	poll.Status = models.PollClosed
	if err := app.Repos.Poll.Update(poll); err != nil {
		return true, fmt.Errorf("update poll: %w", err)
	}

	app.Logger.Info("user closed poll",
		zap.String("userId", post.UserId),
		zap.String("pollId", strconv.Itoa(id)))

	err = RespondToPost(ctx, post, fmt.Sprintf("Poll %s closed", poll.Question))
	return true, err
}
