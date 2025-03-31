package handler

import (
	"context"
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/application"
	"github.com/iamnalinor/vk-mattermost-polls/internal/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

var pendingPolls = map[string]*models.Poll{}

func buildKey(post *model.Post) string {
	return fmt.Sprintf("%s_%s", post.ChannelId, post.UserId)
}

func newPollHandler(ctx context.Context, post *model.Post) (bool, error) {
	app := ctx.Value("app").(*application.Application)
	app.Logger.Debug("user started poll creation", zap.String("channelId", post.ChannelId), zap.String("userId", post.UserId))
	pendingPolls[buildKey(post)] = &models.Poll{}

	err := RespondToPost(ctx, post, "Great! Enter the name of poll. If you want to cancel, type `!cancel`")
	return true, err
}

func pollTextHandler(ctx context.Context, post *model.Post) (bool, error) {
	if strings.HasPrefix(post.Message, "!") {
		return false, nil
	}
	if _, ok := pendingPolls[buildKey(post)]; !ok {
		return false, nil
	}

	app := ctx.Value("app").(*application.Application)
	var message string

	poll := pendingPolls[buildKey(post)]
	if poll.Question == "" {
		app.Logger.Debug("user set question",
			zap.String("channelId", post.ChannelId),
			zap.String("userId", post.UserId),
			zap.String("question", post.Message))
		poll.Question = post.Message
		message = "Enter the first option of the poll:"
	} else {
		app.Logger.Debug("user added option",
			zap.String("channelId", post.ChannelId),
			zap.String("userId", post.UserId),
			zap.String("option", post.Message))
		poll.Options = append(poll.Options, post.Message)
		message = fmt.Sprintf("Options added: %d. Type `!done` once you're done, or add more options", len(poll.Options))
	}
	err := RespondToPost(ctx, post, message)
	return true, err
}

func pollDoneHandler(ctx context.Context, post *model.Post) (bool, error) {
	if _, ok := pendingPolls[buildKey(post)]; !ok {
		return false, nil
	}
	poll := pendingPolls[buildKey(post)]
	if len(poll.Options) == 0 {
		err := RespondToPost(ctx, post, "You need to add at least one option")
		return true, err
	}

	delete(pendingPolls, buildKey(post))
	poll.Status = models.PollActive
	poll.CreatorID = post.UserId

	pollRepo := ctx.Value("app").(*application.Application).Repos.Poll

	id, err := pollRepo.Create(*poll)
	if err != nil {
		return true, fmt.Errorf("create poll: %w", err)
	}

	app := ctx.Value("app").(*application.Application)
	app.Logger.Info("user created poll",
		zap.String("userId", post.UserId),
		zap.String("pollId", strconv.Itoa(id)))

	err = RespondToPost(ctx, post, fmt.Sprintf("Poll created! View it via `!poll %d`", id))
	return true, err
}

func pollCancelHandler(ctx context.Context, post *model.Post) (bool, error) {
	if _, ok := pendingPolls[buildKey(post)]; !ok {
		return false, nil
	}

	app := ctx.Value("app").(*application.Application)
	app.Logger.Debug("user cancelled poll creation",
		zap.String("channelId", post.ChannelId),
		zap.String("userId", post.UserId))

	delete(pendingPolls, buildKey(post))
	err := RespondToPost(ctx, post, "Cancelled.")
	return true, err
}
