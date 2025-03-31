package polling

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/application"
	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
	"strings"
)

type Handler func(ctx context.Context, event *model.Post) (handled bool, err error)

type Dispatcher struct {
	prefixes []string
	handlers []Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		prefixes: make([]string, 0),
		handlers: make([]Handler, 0),
	}
}

func (d *Dispatcher) Register(prefix string, handler Handler) {
	d.prefixes = append(d.prefixes, prefix)
	d.handlers = append(d.handlers, handler)
}

func (d *Dispatcher) Dispatch(ctx context.Context, event *model.WebSocketEvent) (handled bool, err error) {
	if event.EventType() != model.WebsocketEventPosted {
		return false, nil
	}

	postString := event.GetData()["post"].(string)
	post := &model.Post{}
	err = json.Unmarshal([]byte(postString), &post)
	if err != nil {
		return false, fmt.Errorf("unmarshal post: %w", err)
	}

	logger := ctx.Value("app").(*application.Application).Logger

	for i, prefix := range d.prefixes {
		if strings.HasPrefix(post.Message, prefix) {
			logger.Debug("dispatched post", zap.String("postId", post.Id), zap.String("prefix", prefix))
			handled, err := d.handlers[i](ctx, post)
			if err != nil {
				return true, fmt.Errorf("handle prefix %s: %w", prefix, err)
			}
			if handled {
				return true, nil
			}
			logger.Debug("post is rejected by handler", zap.String("postId", post.Id), zap.String("prefix", prefix))
		}
	}

	logger.Debug("post is not handled", zap.String("postId", post.Id))

	return false, nil
}
