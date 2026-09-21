package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourcesByTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResourcesByTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourcesByTagLogic {
	return &ResourcesByTagLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ResourcesByTagLogic) ResourcesByTag(req *types.ResourcesByTagRequest) (resp *types.ResourcesByTagResponse, err error) {
	resp = new(types.ResourcesByTagResponse)

	rpcResp, err := l.svcCtx.TagRPC.ResourcesByTag(l.ctx, &tag.ResourcesByTagRequest{
		TagId:    req.TagId,
		BizId:    req.BizId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[ResourcesByTag] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.ResourceItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		items = append(items, &types.ResourceItem{
			TargetId:   item.TargetId,
			BizId:      item.BizId,
			CreateTime: item.CreateTime,
		})
	}
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
