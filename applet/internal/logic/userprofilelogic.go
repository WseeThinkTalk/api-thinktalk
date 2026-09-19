package logic

import (
	"context"

	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
	"api-thinktalk/client/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileLogic) UserProfile(req *types.UserProfileRequest) (resp *types.UserProfileResponse, err error) {
	id, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: req.UserId})
	if err != nil {
		logx.Errorf("UserProfile findById error: %v", err)
		return nil, err
	}

	return &types.UserProfileResponse{
		UserId:   id.UserId,
		Username: id.Username,
		Avatar:   id.Avatar,
		Bio:      id.Bio,
	}, nil
}
