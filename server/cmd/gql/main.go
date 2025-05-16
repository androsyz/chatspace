package main

import (
	"context"
	"hora-server/cmd/initalize"
	"hora-server/config"
	"hora-server/graph/generated"
	"hora-server/handler/middleware"
	"hora-server/handler/resolver"
	"log"

	http "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	var ctx = context.Background()
	zlog := config.NewLogger()

	cfg, err := config.NewViper()
	if err != nil {
		zlog.Err(err)
		return
	}

	app, err := initalize.Bootstrap(ctx, cfg, zlog)
	if err != nil {
		zlog.Err(err)
		return
	}

	rsvl, err := resolver.NewResolver(app.UcUser)
	if err != nil {
		zlog.Err(err)
		return
	}

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: rsvl}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	address := cfg.Server.Address
	if address == "" {
		address = defaultPort
	}

	httpServer := http.NewServer(
		http.Address(":"+address),
		http.Middleware(
			middleware.AuthenticationGQL(cfg.Settings.JWTSecret),
		),
	)

	httpServer.Handle("/", playground.Handler("GraphQL playground", "/query"))
	httpServer.Handle("/query", srv)

	zlog.Info().Msgf("connect to http://localhost:%s for GraphQL playground", address)
	log.Fatal(httpServer.Start(ctx))
}
