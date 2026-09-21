// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	user "api-thinktalk/client/user/service"
	"context"
	"encoding/json"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	resp = new(types.UserInfoResponse)

	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return
	}

	id, err := l.svcCtx.UserRPC.FindById(l.ctx,
		&user.FindByIdRequest{UserId: userId})
	if err != nil {
		logx.Errorf("findById error: %v", err)
		return nil, err
	}

	resp.UserId = id.UserId
	resp.Username = id.Username
	resp.Avatar = id.Avatar
	resp.Role = id.Role
	resp.DisplayId = id.DisplayId
	resp.Bio = id.Bio
	resp.Gender = id.Gender
	resp.ProfileCover = id.ProfileCover
	return resp, nil
}
