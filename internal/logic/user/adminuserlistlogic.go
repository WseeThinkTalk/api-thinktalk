package user

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	user "api-thinktalk/client/user/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserListLogic {
	return &AdminUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserListLogic) AdminUserList(req *types.AdminUserListRequest) (resp *types.AdminUserListResponse, err error) {
	resp = new(types.AdminUserListResponse)

	rpcResp, err := l.svcCtx.UserRPC.AdminUserList(l.ctx, &user.AdminUserListRequest{
		Keyword:  req.Keyword,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var users []types.AdminUserInfo
	if rpcResp != nil && rpcResp.Data != nil {
		for _, v := range rpcResp.Data.Users {
			if v != nil {
				users = append(users, types.AdminUserInfo{
					UserId:       v.UserId,
					Username:     v.Username,
					Avatar:       v.Avatar,
					Role:         v.Role,
					DisplayId:    v.DisplayId,
					Bio:          v.Bio,
					Gender:       v.Gender,
					ProfileCover: v.ProfileCover,
				})
			}
		}
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}

	resp.Users = users
	return resp, nil
}
