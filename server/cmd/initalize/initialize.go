package initalize

import (
	"context"
	"hora-server/config"
	"hora-server/repository"
	"hora-server/usecase"

	"github.com/rs/zerolog"
)

type App struct {
	UcUser *usecase.UcUser
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

	// setup repository
	zlog.Info().Msg("Initialize Repository")
	repoUser := repository.NewUserRepository(dbConn)

	// setup usecase
	zlog.Info().Msg("Initialize Usecase")
	ucUser := usecase.NewUserUsecase(cfg, repoUser)

	return App{
		UcUser: ucUser,
	}, nil
}
