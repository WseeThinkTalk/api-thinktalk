package social

import (
	"context"
	"encoding/json"

	concernedpb "api-thinktalk/client/concerned/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedCountLogic {
	return &ConcernedCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConcernedCountLogic) ConcernedCount(req *types.ConcernedCountRequest) (resp *types.ConcernedCountResponse, err error) {
	resp = new(types.ConcernedCountResponse)

	objId := req.ObjId
	if objId == 0 {
		userIdVal := l.ctx.Value("userId")
		if userIdVal != nil {
			if num, ok := userIdVal.(json.Number); ok {
				uid, _ := num.Int64()
				objId = -uid
			}
		}
	}

	rpcResp, err := l.svcCtx.ConcernedRPC.ConcernedCount(l.ctx, &concernedpb.ConcernedCountRequest{
		BizId: req.BizId,
		ObjId: objId,
	})
	if err != nil {
		l.Errorf("[ConcernedCount] rpc err: %v", err)
		return nil, err
	}
	if rpcResp != nil && rpcResp.Data != nil {
		resp.Count = rpcResp.Data.ConcernedNum
	}
	return resp, nil
}
