package social

import (
	"context"

	"api-thinktalk/client/message/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadCountLogic {
	return &UnreadCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnreadCountLogic) UnreadCount(userId int64) (resp *types.UnreadCountResponse, err error) {
	resp = new(types.UnreadCountResponse)

	rpcResp, err := l.svcCtx.MessageRPC.UnreadCount(l.ctx, &pb.UnreadCountRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[UnreadCount] rpc err: %v", err)
		return nil, err
	}
	resp.Total = rpcResp.Total
	resp.TypeCounts = rpcResp.TypeCounts
	return resp, nil
}
