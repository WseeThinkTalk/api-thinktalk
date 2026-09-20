package social

import (
	"context"

	"api-thinktalk/client/follow/pb"
	user "api-thinktalk/client/user/service"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FansListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansListLogic {
	return &FansListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FansListLogic) FansList(userId int64, req *types.FansListRequest) (*types.FansListResponse, error) {
	resp, err := l.svcCtx.FollowRPC.FansList(l.ctx, &pb.FansListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[FansList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.FollowItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		targetName := ""
		targetAvatar := ""
		if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: item.FansUserId}); err == nil {
			targetName = userResp.Username
			targetAvatar = userResp.Avatar
		}
		items = append(items, &types.FollowItem{
			Id:               item.UserId,
			FollowedUserId:   item.FansUserId,
			CreateTime:       item.CreateTime,
			FansCount:        item.FansCount,
			TargetUserName:   targetName,
			TargetUserAvatar: targetAvatar,
		})
	}
	return &types.FansListResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
