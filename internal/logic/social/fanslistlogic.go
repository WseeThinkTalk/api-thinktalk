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

func (l *FansListLogic) FansList(userId int64, req *types.FansListRequest) (resp *types.FansListResponse, err error) {
	resp = new(types.FansListResponse)
	resp.Items = make([]*types.FansItem, 0)

	rpcResp, err := l.svcCtx.FollowRPC.FansList(l.ctx, &pb.FansListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[FansList] rpc err: %v", err)
		return nil, err
	}

	// 转换粉丝列表项并补充对方用户信息
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.FansItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			targetName := ""
			targetAvatar := ""
			if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: v.FansUserId}); err == nil && userResp != nil && userResp.Data != nil {
				targetName = userResp.Data.Username
				targetAvatar = userResp.Data.Avatar
			}
			items = append(items, &types.FansItem{
				Id:             v.UserId,
				FansId:         v.FansUserId,
				FansName:       targetName,
				FansAvatar:     targetAvatar,
				CreateTime:     v.CreateTime,
			})
		}
		resp.Items = items
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
