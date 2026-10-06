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

	if id != nil && id.Data != nil {
		resp.UserId = id.Data.UserId
		resp.Username = id.Data.Username
		resp.Avatar = id.Data.Avatar
		resp.Role = id.Data.Role
		resp.DisplayId = id.Data.DisplayId
		resp.Bio = id.Data.Bio
		resp.Gender = id.Data.Gender
		resp.ProfileCover = id.Data.ProfileCover
	}
	return resp, nil
}
