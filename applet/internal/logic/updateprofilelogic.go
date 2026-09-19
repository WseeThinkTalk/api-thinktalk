// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
	"api-thinktalk/client/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileRequest) (resp *types.UpdateProfileResponse, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return nil, fmt.Errorf("user not login")
	}

	_, err = l.svcCtx.UserRPC.UpdateProfile(l.ctx, &user.UpdateProfileRequest{
		UserId:       userId,
		Username:     req.Username,
		Avatar:       req.Avatar,
		Bio:          req.Bio,
		Gender:       req.Gender,
		ProfileCover: req.ProfileCover,
	})
	if err != nil {
		logx.Errorf("UpdateProfile error: %v", err)
		return nil, err
	}

	return &types.UpdateProfileResponse{}, nil
}
