package user

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	user "api-thinktalk/client/user/service"

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
	resp = new(types.UserProfileResponse)

	id, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: req.UserId})
	if err != nil {
		logx.Errorf("UserProfile findById error: %v", err)
		return nil, err
	}

	resp.UserId = id.UserId
	resp.Username = id.Username
	resp.Avatar = id.Avatar
	resp.Bio = id.Bio
	return resp, nil
}
