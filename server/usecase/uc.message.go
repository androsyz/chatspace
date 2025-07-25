package usecase

import (
	"context"
	"encoding/json"
	"chatspace-server/constant"
	"chatspace-server/graph/model"
	modelDB "chatspace-server/model"
	"chatspace-server/pkg/authctx"
	"chatspace-server/pkg/gqlhelper"
	"chatspace-server/pkg/helper"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type repoMessageInterface interface {
	Create(ctx context.Context, message *modelDB.MessageDB) (*string, error)
	GetMessages(ctx context.Context) ([]*modelDB.MessageDB, error)
	PublishMessage(ctx context.Context, spaceID string, data []byte) error
	SubscribeMessage(ctx context.Context, spaceID string) *redis.PubSub
}

type UcMessage struct {
	repoMessage repoMessageInterface
	repoUser    repoUserInterface
	repoSpace   repoSpaceInterface
	zlog        zerolog.Logger
}

func NewMessageUseCase(repoMessage repoMessageInterface, repoUser repoUserInterface, repoSpace repoSpaceInterface, zlog zerolog.Logger) *UcMessage {
	return &UcMessage{
		repoMessage: repoMessage,
		repoUser:    repoUser,
		repoSpace:   repoSpace,
		zlog:        zlog,
	}
}

func (uc *UcMessage) SendMessage(ctx context.Context, spaceID string, content string) (*model.Message, error) {
	userID, err := authctx.GetAuthUserID(ctx)
	if err != nil {
		return nil, err
	}

	userUUID, err := helper.StrToUUID(userID)
	if err != nil {
		return nil, err
	}

	spaceUUID, err := helper.StrToUUID(spaceID)
	if err != nil {
		return nil, err
	}

	payload := &modelDB.MessageDB{
		Content: content,
		UserID:  *userUUID,
		SpaceID: *spaceUUID,
	}

	msgID, err := uc.repoMessage.Create(ctx, payload)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrCreatingField("message"), err)
	}

	resp := &model.Message{
		ID:      *msgID,
		Content: content,
		User:    &model.User{ID: userUUID.String()},
		Space:   &model.Space{ID: spaceUUID.String()},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		uc.zlog.Error().Err(err).Msg(constant.ErrMsgMarshal)
		return nil, err
	}

	err = uc.repoMessage.PublishMessage(ctx, spaceID, data)
	if err != nil {
		uc.zlog.Error().Err(err).Msg(constant.ErrMsgPublish)
		return nil, err
	}

	return resp, nil
}

func (uc *UcMessage) Messages(ctx context.Context, spaceID string) ([]*model.Message, error) {
	_, err := authctx.GetAuthUserID(ctx)
	if err != nil {
		return nil, err
	}

	messages, err := uc.repoMessage.GetMessages(ctx)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrGetField("messages"), err)
	}

	var resp []*model.Message
	for _, m := range messages {
		temp, err := uc.PopulateMessageField(ctx, m)
		if err != nil {
			uc.zlog.Error().Err(err)
		}

		resp = append(resp, temp)
	}

	return resp, nil
}

func (uc *UcMessage) MessageSent(ctx context.Context, spaceID string) (<-chan *model.Message, error) {
	ch := make(chan *model.Message, 1)

	pubsub := uc.repoMessage.SubscribeMessage(ctx, spaceID)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		uc.zlog.Error().Err(err).Msg(constant.ErrMsgSubscribe)
		_ = pubsub.Close()
		close(ch)
		return ch, err
	}

	chRedis := pubsub.Channel()

	go func() {
		defer func() {
			_ = pubsub.Close()
			close(ch)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-chRedis:
				if !ok {
					return
				}

				var message model.Message
				err := json.Unmarshal([]byte(msg.Payload), &message)
				if err != nil {
					uc.zlog.Error().Err(err).Msg(constant.ErrMsgUnmarshal)
					continue
				}

				select {
				case ch <- &message:
				default:
					uc.zlog.Warn().Msg(constant.ErrMsgSubsFull)
				}
			}
		}
	}()

	return ch, nil
}

func (uc *UcMessage) PopulateMessageField(ctx context.Context, message *modelDB.MessageDB) (*model.Message, error) {
	resp := &model.Message{
		ID:      message.ID.String(),
		Content: message.Content,
	}

	if gqlhelper.IsCalled(ctx, "user") {
		user, err := uc.repoUser.GetByID(ctx, message.UserID.String())
		if err != nil {
			return nil, constant.ErrUserNotFound
		}

		tempUser := &model.User{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		}

		resp.User = tempUser
	}

	if gqlhelper.IsCalled(ctx, "space") {
		space, err := uc.repoSpace.GetSpaceByID(ctx, message.SpaceID.String())
		if err != nil {
			return nil, constant.ErrWithMsg(constant.ErrGetField("space"), err)
		}

		tempSpace := &model.Space{
			ID:          space.ID.String(),
			Name:        space.Name,
			Description: &space.Description,
		}

		resp.Space = tempSpace
	}

	return resp, nil
}
