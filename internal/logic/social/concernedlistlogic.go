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

func (l *ConcernedListLogic) List(userId int64, req *types.ConcernedListRequest) (*types.ConcernedListResponse, error) {
	resp, err := l.svcCtx.ConcernedRPC.ConcernedList(l.ctx, &concernedpb.ConcernedListRequest{
		UserId:   userId,
		BizId:    req.BizId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[ConcernedList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.ConcernedItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.ConcernedItem{
			Id:         item.Id,
			BizId:      item.BizId,
			ObjId:      item.ObjId,
			CreateTime: item.CreateTime,
		})
	}
	return &types.ConcernedListResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
