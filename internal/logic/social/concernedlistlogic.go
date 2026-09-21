package social

import (
	"context"

	concernedpb "api-thinktalk/client/concerned/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedListLogic {
	return &ConcernedListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConcernedListLogic) List(userId int64, req *types.ConcernedListRequest) (resp *types.ConcernedListResponse, err error) {
	resp = new(types.ConcernedListResponse)
	resp.Items = make([]*types.ConcernedItem, 0)

	rpcResp, err := l.svcCtx.ConcernedRPC.ConcernedList(l.ctx, &concernedpb.ConcernedListRequest{
		UserId:   userId,
		BizId:    req.BizId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[ConcernedList] rpc err: %v", err)
		return nil, err
	}

	// 转换关注动态记录项
	items := make([]*types.ConcernedItem, 0, len(rpcResp.Items))
	for _, v := range rpcResp.Items {
		items = append(items, &types.ConcernedItem{
			Id:         v.Id,
			BizId:      v.BizId,
			ObjId:      v.ObjId,
			CreateTime: v.CreateTime,
		})
	}
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
