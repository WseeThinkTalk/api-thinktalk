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

	var items []*types.AdminUserItem
	for _, item := range rpcResp.Items {
		items = append(items, &types.AdminUserItem{
			UserId:       item.UserId,
			Username:     item.Username,
			Mobile:       item.Mobile,
			Avatar:       item.Avatar,
			Role:         item.Role,
			DisplayId:    item.DisplayId,
			Bio:          item.Bio,
			Gender:       item.Gender,
			ProfileCover: item.ProfileCover,
			CreateTime:   item.CreateTime,
		})
	}

	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
