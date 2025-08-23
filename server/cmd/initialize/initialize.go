package initialize

import (
	"chatspace-server/config"
	"chatspace-server/repository"
	"chatspace-server/usecase"
	"context"

	"github.com/rs/zerolog"
)

type App struct {
	UcUser    *usecase.UcUser
	UcSpace   *usecase.UcSpace
	UcMessage *usecase.UcMessage
}

func Bootstrap(ctx context.Context, cfg *config.Config, zlog zerolog.Logger) (App, error) {
	app := App{}

	// setup database
	zlog.Info().Msg("Initialize Database")
	dbConn, err := config.NewDatabase(cfg.Database)
	if err != nil {
		zlog.Error().Err(err).Msg("Failed initialize database")
		return app, err
	}

	// setup redis
	zlog.Info().Msg("Initialize Redis")
	rdsConn, err := config.NewRedis(cfg.Redis)
	if err != nil {
		zlog.Error().Err(err).Msg("Failed initialize redis")
		return app, err
	}

	// setup repository
	zlog.Info().Msg("Initialize Repository")
	repoUser := repository.NewUserRepository(dbConn)
	repoSpace := repository.NewSpaceRepository(dbConn)
	repoMessage := repository.NewMessageRepository(dbConn, rdsConn)

	// setup usecase
	zlog.Info().Msg("Initialize Usecase")
	ucUser := usecase.NewUserUsecase(cfg, repoUser, zlog)
	ucSpace := usecase.NewSpaceUseCase(repoSpace, zlog)
	ucMessage := usecase.NewMessageUseCase(repoMessage, repoUser, repoSpace, zlog)

	return App{
		UcUser:    ucUser,
		UcSpace:   ucSpace,
		UcMessage: ucMessage,
	}, nil
}
