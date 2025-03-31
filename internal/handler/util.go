package handler

import (
	"context"
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/application"
	"github.com/mattermost/mattermost-server/v6/model"
)

func RespondToPost(ctx context.Context, post *model.Post, message string) error {
	mmClient := ctx.Value("app").(*application.Application).MmClient

	_, _, err := mmClient.CreatePost(&model.Post{
		ChannelId: post.ChannelId,
		Message:   message,
		RootId:    post.RootId,
	})
	if err != nil {
		return fmt.Errorf("respond to post: %w", err)
	}
	return nil
}
