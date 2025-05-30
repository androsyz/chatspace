package usecase

import (
	"context"
	"hora-server/constant"
	"hora-server/graph/model"
	modelDB "hora-server/model"
	"hora-server/pkg/authctx"
	"hora-server/pkg/gqlhelper"
	"hora-server/pkg/helper"

	"github.com/rs/zerolog"
)

type repoSpaceInterface interface {
	Create(ctx context.Context, space *modelDB.SpaceDB) (*string, error)
	CreateSpaceMember(ctx context.Context, spaceMember *modelDB.SpaceMemberDB) error
	GetSpaceByID(ctx context.Context, id string) (*modelDB.SpaceDB, error)
	GetSpaceMember(ctx context.Context, spaceID string) ([]*modelDB.SpaceMemberDB, error)
	GetMemberBySpaceID(ctx context.Context, spaceID, role string) ([]*modelDB.UserDB, error)
	GetSpaces(ctx context.Context) ([]*modelDB.SpaceDB, error)
}

type UcSpace struct {
	repoSpace repoSpaceInterface
	zlog      zerolog.Logger
}

func NewSpaceUseCase(repoSpace repoSpaceInterface, zlog zerolog.Logger) *UcSpace {
	return &UcSpace{
		repoSpace: repoSpace,
		zlog:      zlog,
	}
}

func (uc *UcSpace) CreateSpace(ctx context.Context, request model.SpaceRequest) (*model.Space, error) {
	userID, err := authctx.GetAuthUserID(ctx)
	if err != nil {
		return nil, err
	}

	if request.Name == "" {
		return nil, constant.ErrMissingField("name")
	}

	payload := &modelDB.SpaceDB{
		Name:        request.Name,
		Description: *request.Description,
	}

	spaceID, err := uc.repoSpace.Create(ctx, payload)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrCreatingField("space"), err)
	}

	userUUID, err := helper.StrToUUID(userID)
	if err != nil {
		return nil, err
	}

	spaceUUID, err := helper.StrToUUID(*spaceID)
	if err != nil {
		return nil, err
	}

	memberPayload := &modelDB.SpaceMemberDB{
		UserID:  *userUUID,
		SpaceID: *spaceUUID,
		Role:    constant.ROLE_ADMIN,
	}

	err = uc.repoSpace.CreateSpaceMember(ctx, memberPayload)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrCreatingField("space member"), err)
	}

	resp := &model.Space{
		ID:   *spaceID,
		Name: request.Name,
	}

	return resp, nil
}

func (uc *UcSpace) JoinSpace(ctx context.Context, spaceID string) (*model.Space, error) {
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

	payload := &modelDB.SpaceMemberDB{
		UserID:  *userUUID,
		SpaceID: *spaceUUID,
		Role:    constant.ROLE_MEMBER,
	}

	err = uc.repoSpace.CreateSpaceMember(ctx, payload)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrCreatingField("space member"), err)
	}

	space, err := uc.repoSpace.GetSpaceByID(ctx, spaceID)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrGetField("space"), err)
	}

	resp := &model.Space{
		ID:          space.ID.String(),
		Name:        space.Name,
		Description: &space.Description,
	}

	return resp, nil
}

func (uc *UcSpace) Spaces(ctx context.Context) ([]*model.Space, error) {
	_, err := authctx.GetAuthUserID(ctx)
	if err != nil {
		return nil, err
	}

	spaces, err := uc.repoSpace.GetSpaces(ctx)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrGetField("spaces"), err)
	}

	var resp []*model.Space
	for _, s := range spaces {
		temp, err := uc.PopulateSpaceField(ctx, *s)
		if err != nil {
			uc.zlog.Error().Err(err)
		}

		resp = append(resp, temp)
	}

	return resp, nil
}

func (uc *UcSpace) Space(ctx context.Context, id string) (*model.Space, error) {
	_, err := authctx.GetAuthUserID(ctx)
	if err != nil {
		return nil, err
	}

	space, err := uc.repoSpace.GetSpaceByID(ctx, id)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrGetField("space"), err)
	}

	resp, err := uc.PopulateSpaceField(ctx, *space)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (uc *UcSpace) PopulateSpaceField(ctx context.Context, space modelDB.SpaceDB) (*model.Space, error) {
	spaceID := space.ID.String()
	resp := &model.Space{
		ID:          spaceID,
		Name:        space.Name,
		Description: &space.Description,
	}

	if gqlhelper.IsCalled(ctx, "members") {
		members, err := uc.repoSpace.GetMemberBySpaceID(ctx, spaceID, constant.ROLE_MEMBER)
		if err != nil {
			return nil, constant.ErrWithMsg(constant.ErrGetField("member"), err)
		}

		var respMembers []*model.User
		for _, m := range members {
			temp := &model.User{
				ID:    m.ID.String(),
				Email: m.Email,
				Name:  m.Name,
			}
			respMembers = append(respMembers, temp)
		}

		resp.Members = respMembers
	}

	if gqlhelper.IsCalled(ctx, "admins") {
		admins, err := uc.repoSpace.GetMemberBySpaceID(ctx, spaceID, constant.ROLE_ADMIN)
		if err != nil {
			return nil, constant.ErrWithMsg(constant.ErrGetField("member"), err)
		}

		var respAdmins []*model.User
		for _, a := range admins {
			temp := &model.User{
				ID:    a.ID.String(),
				Email: a.Email,
				Name:  a.Name,
			}
			respAdmins = append(respAdmins, temp)
		}

		resp.Admins = respAdmins
	}

	return resp, nil
}
