package application

import (
	"context"
	"fmt"
	"github.com/iamnalinor/vk-mattermost-polls/internal/config"
	"github.com/iamnalinor/vk-mattermost-polls/internal/repo"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
	"time"
)

type Application struct {
	Conn     *tarantool.Connection
	Repos    *repo.Repositories
	Logger   *zap.Logger
	MmClient *model.Client4
}

func NewApplication(ctx context.Context) (*Application, error) {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("init zap logger: %w", err)
	}

	conn, err := newTarantoolConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to tarantool: %w", err)
	}
	logger.Info("connected to tarantool")

	repos, err := repo.NewRepositories(conn)
	if err != nil {
		return nil, fmt.Errorf("init repositories: %w", err)
	}

	mmClient, err := newMattermostClient(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("init mattermost client: %w", err)
	}

	return &Application{
		Conn:     conn,
		Repos:    repos,
		Logger:   logger,
		MmClient: mmClient,
	}, nil
}

func newTarantoolConnection(ctx context.Context) (*tarantool.Connection, error) {
	cfg := ctx.Value("config").(config.Config)
	dialer := tarantool.NetDialer{
		Address:  cfg.TarantoolAddress,
		User:     cfg.TarantoolUser,
		Password: cfg.TarantoolPassword,
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newMattermostClient(ctx context.Context, logger *zap.Logger) (*model.Client4, error) {
	cfg := ctx.Value("config").(config.Config)
	client := model.NewAPIv4Client(cfg.MattermostUrl)
	client.SetToken(cfg.MattermostToken)

	user, _, err := client.GetUser("me", "")
	if err != nil {
		return nil, fmt.Errorf("log in to mattermost: %w", err)
	}
	logger.Info("logged in to mattermost", zap.String("userId", user.Id), zap.String("username", user.Username))

	return client, nil
}
