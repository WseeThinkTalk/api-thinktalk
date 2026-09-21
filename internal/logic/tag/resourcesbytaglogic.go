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
	// 转换标签关联资源项
	items := make([]*types.ResourceItem, 0, len(rpcResp.Items))
	for _, v := range rpcResp.Items {
		items = append(items, &types.ResourceItem{
			TargetId:   v.TargetId,
			BizId:      v.BizId,
			CreateTime: v.CreateTime,
		})
	}
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
