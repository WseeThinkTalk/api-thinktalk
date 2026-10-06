package social

import (
	"context"

	"api-thinktalk/client/follow/pb"
	user "api-thinktalk/client/user/service"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FollowListLogic) FollowList(userId int64, req *types.FollowListRequest) (resp *types.FollowListResponse, err error) {
	resp = new(types.FollowListResponse)
	resp.Items = make([]*types.FollowItem, 0)

	rpcResp, err := l.svcCtx.FollowRPC.FollowList(l.ctx, &pb.FollowListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[FollowList] rpc err: %v", err)
		return nil, err
	}

	// 转换关注列表项并补充对方用户信息
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.FollowItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			targetName := ""
			targetAvatar := ""
			if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: v.FollowedUserId}); err == nil && userResp != nil && userResp.Data != nil {
				targetName = userResp.Data.Username
				targetAvatar = userResp.Data.Avatar
			}
			items = append(items, &types.FollowItem{
				Id:           v.Id,
				FollowId:     v.FollowedUserId,
				FollowName:   targetName,
				FollowAvatar: targetAvatar,
				CreateTime:   v.CreateTime,
			})
		}
		resp.Items = items
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
