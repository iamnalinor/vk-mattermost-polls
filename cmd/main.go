package main

import (
	"context"
	"github.com/iamnalinor/vk-mattermost-polls/internal/application"
	"github.com/iamnalinor/vk-mattermost-polls/internal/config"
	"github.com/iamnalinor/vk-mattermost-polls/internal/handler"
	"github.com/iamnalinor/vk-mattermost-polls/internal/polling"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "config", cfg))
	defer cancel()

	app, err := application.NewApplication(ctx)
	if err != nil {
		log.Fatalf("init application: %v", err)
	}

	dispatcher := polling.NewDispatcher()
	handler.RegisterHandlers(dispatcher)

	ctx = context.WithValue(context.WithValue(ctx, "app", app), "dispatcher", dispatcher)

	go func() {
		err := polling.Polling(ctx)
		if err != nil {
			log.Fatalf("polling: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit

	app.Logger.Info("received signal", zap.String("signal", sig.String()))
	cancel()
}
