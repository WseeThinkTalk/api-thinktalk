package social

import (
	"context"

	"api-thinktalk/client/message/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MarkReadLogic) MarkRead(userId int64, req *types.MarkReadRequest) (*types.MarkReadResponse, error) {
	_, err := l.svcCtx.MessageRPC.MarkRead(l.ctx, &pb.MarkReadRequest{
		UserId:         userId,
		NotificationId: req.NotificationId,
	})
	if err != nil {
		l.Errorf("[MarkRead] rpc err: %v", err)
		return nil, err
	}
	return &types.MarkReadResponse{}, nil
}
