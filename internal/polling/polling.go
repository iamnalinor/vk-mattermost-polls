package polling

import (
	"context"
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/application"
	"github.com/iamnalinor/vk-mattermost-polls/internal/config"
	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
	"net/url"
	"time"
)

func Polling(ctx context.Context) error {
	cfg := ctx.Value("config").(config.Config)
	app := ctx.Value("app").(*application.Application)

	mattermostUrl, err := url.Parse(cfg.MattermostUrl)
	if err != nil {
		return fmt.Errorf("parse mattermost url: %w", err)
	}

	for {
		app.Logger.Debug("connecting to mattermost websocket...")
		wsClient, err := model.NewWebSocketClient4(
			fmt.Sprintf("ws://%s", mattermostUrl.Host+mattermostUrl.Path),
			cfg.MattermostToken,
		)
		if err != nil {
			app.Logger.Error("connect to mattermost websocket", zap.Error(err))

			<-time.After(time.Second)
			continue
		}
		app.Logger.Debug("connected to mattermost websocket")

		err = listenToEvents(ctx, wsClient)
		app.Logger.Error("listen to events", zap.Error(err))

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func listenToEvents(ctx context.Context, wsClient *model.WebSocketClient) error {
	wsClient.Listen()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-wsClient.EventChannel:
			go handleEvent(ctx, event)
		}
		if wsClient.ListenError != nil {
			return wsClient.ListenError
		}
	}
}

func handleEvent(ctx context.Context, event *model.WebSocketEvent) {
	logger := ctx.Value("app").(*application.Application).Logger
	logger.Debug("got event", zap.String("event", fmt.Sprint(event.GetData())))

	dispatcher := ctx.Value("dispatcher").(*Dispatcher)
	_, err := dispatcher.Dispatch(ctx, event)
	if err != nil {
		logger.Error("dispatch event", zap.Error(err))
	}
}
