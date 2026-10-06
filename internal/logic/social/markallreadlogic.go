package social

import (
	"context"

	"api-thinktalk/client/message/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkAllReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllReadLogic {
	return &MarkAllReadLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MarkAllReadLogic) MarkAllRead(userId int64, req *types.MarkAllReadRequest) (resp *types.MarkAllReadResponse, err error) {
	resp = new(types.MarkAllReadResponse)

	_, err = l.svcCtx.MessageRPC.MarkAllRead(l.ctx, &pb.MarkAllReadRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[MarkAllRead] rpc err: %v", err)
		return nil, err
	}
	return resp, nil
}
