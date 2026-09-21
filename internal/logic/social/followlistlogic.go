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

	items := make([]*types.FollowItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		targetName := ""
		targetAvatar := ""
		if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: item.FollowedUserId}); err == nil {
			targetName = userResp.Username
			targetAvatar = userResp.Avatar
		}
		items = append(items, &types.FollowItem{
			Id:               item.Id,
			FollowedUserId:   item.FollowedUserId,
			CreateTime:       item.CreateTime,
			FansCount:        item.FansCount,
			TargetUserName:   targetName,
			TargetUserAvatar: targetAvatar,
		})
	}
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
