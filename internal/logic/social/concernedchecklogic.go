package social

import (
	"context"

	concernedpb "api-thinktalk/client/concerned/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedCheckLogic {
	return &ConcernedCheckLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConcernedCheckLogic) Check(userId int64, req *types.ConcernedCheckRequest) (*types.ConcernedCheckResponse, error) {
	resp, err := l.svcCtx.ConcernedRPC.IsConcerned(l.ctx, &concernedpb.IsConcernedRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[ConcernedCheck] rpc err: %v", err)
		return nil, err
	}
	return &types.ConcernedCheckResponse{IsConcerned: resp.IsConcerned}, nil
}
