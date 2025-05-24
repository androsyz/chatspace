package resolver

import (
	"context"
	"hora-server/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ucUserInterface interface {
	Register(ctx context.Context, request model.RegisterRequest) (*model.AuthPayload, error)
	Login(ctx context.Context, request model.LoginRequest) (*model.AuthPayload, error)
	RefreshToken(ctx context.Context, request model.RefreshRequest) (*model.AuthPayload, error)
	Me(ctx context.Context) (*model.User, error)
}

func NewResolver(
	ucUser ucUserInterface,
) (*Resolver, error) {
	return &Resolver{
		ucUser: ucUser,
	}, nil
}

type Resolver struct {
	ucUser ucUserInterface
}
