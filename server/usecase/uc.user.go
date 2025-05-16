package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hora-server/config"
	"hora-server/graph/model"
	modelDB "hora-server/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = fmt.Errorf("user not found")

type repoUserInterface interface {
	Create(ctx context.Context, user *modelDB.UserDB) error
	GetByID(ctx context.Context, id string) (*modelDB.UserDB, error)
	GetByEmail(ctx context.Context, email string) (*modelDB.UserDB, error)
	Update(ctx context.Context, user *modelDB.UserDB) error
	Delete(ctx context.Context, id string) error
}

type UcUser struct {
	cfg      *config.Config
	repoUser repoUserInterface
}

func NewUserUsecase(cfg *config.Config, repoUser repoUserInterface) *UcUser {
	return &UcUser{
		cfg:      cfg,
		repoUser: repoUser,
	}
}

func (uc *UcUser) Register(ctx context.Context, request model.RegisterRequest) (*model.AuthPayload, error) {
	if request.Username == "" || request.Email == "" || request.Password == "" {
		return nil, fmt.Errorf("username, email, and password are required")
	}

	existingUser, err := uc.repoUser.GetByEmail(ctx, request.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to check for existing user: %w", err)
		}
	}
	if existingUser != nil {
		return nil, fmt.Errorf("email is already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	payload := &modelDB.UserDB{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Name:     request.Name,
	}

	err = uc.repoUser.Create(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user := &model.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Name:     request.Name,
	}

	return uc.generateAuthPayload(user)
}

func (uc *UcUser) Login(ctx context.Context, request model.LoginRequest) (*model.AuthPayload, error) {
	if request.Email == "" || request.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := uc.repoUser.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	resp := &model.User{
		Username: user.Username,
	}

	return uc.generateAuthPayload(resp)
}

func (uc *UcUser) RefreshToken(ctx context.Context, request model.RefreshRequest) (*model.AuthPayload, error) {
	token, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.cfg.Settings.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid subject in token")
	}

	user, err := uc.repoUser.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	resp := &model.User{
		Username: user.Username,
	}

	return uc.generateAuthPayload(resp)
}

func (uc *UcUser) Logout(ctx context.Context) (*model.LogoutResponse, error) {
	msg := "logged out"
	return &model.LogoutResponse{
		Message: &msg,
	}, nil
}

func (uc *UcUser) Me(ctx context.Context) (*model.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (uc *UcUser) generateAuthPayload(user *model.User) (*model.AuthPayload, error) {
	accessToken, err := uc.generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.generateRefreshJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &model.AuthPayload{
		Token:        accessToken,
		RefreshToken: &refreshToken,
		User: &model.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Name:     user.Username,
		},
	}, nil
}

func (uc *UcUser) generateJWT(userID string) (string, error) {
	duration := time.Hour * time.Duration(uc.cfg.Settings.TokenDuration)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.cfg.Settings.JWTSecret))
}

func (uc *UcUser) generateRefreshJWT(userID string) (string, error) {
	duration := time.Hour * time.Duration(uc.cfg.Settings.RefreshTokenDuration)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.cfg.Settings.JWTSecret))
}
