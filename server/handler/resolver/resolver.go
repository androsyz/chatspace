package resolver

import (
	"context"
	"hora-server/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ucUserInterface interface {
	Register(ctx context.Context, request model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, request model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, request model.RefreshRequest) (*model.AuthResponse, error)
	User(ctx context.Context) (*model.User, error)
}

type ucSpaceInterface interface {
	CreateSpace(ctx context.Context, request model.SpaceRequest) (*model.Space, error)
	JoinSpace(ctx context.Context, spaceID string) (*model.Space, error)
	Spaces(ctx context.Context) ([]*model.Space, error)
	Space(ctx context.Context, id string) (*model.Space, error)
}

type ucMessageInterface interface {
	SendMessage(ctx context.Context, spaceID string, content string) (*model.Message, error)
	Messages(ctx context.Context, spaceID string) ([]*model.Message, error)
	MessageSent(ctx context.Context, spaceID string) (<-chan *model.Message, error)
}

func NewResolver(
	ucUser ucUserInterface,
	ucSpace ucSpaceInterface,
	ucMessage ucMessageInterface,
) (*Resolver, error) {
	return &Resolver{
		ucUser:    ucUser,
		ucSpace:   ucSpace,
		ucMessage: ucMessage,
	}, nil
}

type Resolver struct {
	ucUser    ucUserInterface
	ucSpace   ucSpaceInterface
	ucMessage ucMessageInterface
}
