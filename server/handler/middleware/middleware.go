package middleware

import (
	"context"
	"strings"

	midKratos "github.com/go-kratos/kratos/v2/middleware"
	trKratos "github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserCtxKey contextKey = "x-user-id"
)

type AuthUser struct {
	UserID     string
	Roles      []string
	Platform   string
	AppVersion string
}

func AuthenticationGQL(secret string) func(midKratos.Handler) midKratos.Handler {
	return func(handler midKratos.Handler) midKratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := trKratos.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return handler(ctx, req)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return nil, err
			}

			userID, _ := claims["sub"].(string)

			authUser := &AuthUser{
				UserID: userID,
			}

			newCtx := context.WithValue(ctx, UserCtxKey, authUser)

			return handler(newCtx, req)
		}
	}
}
