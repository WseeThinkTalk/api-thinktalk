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
		BizId:    req.ResourceType,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[ResourcesByTag] rpc err: %v", err)
		return nil, err
	}
	// 转换标签关联资源项
	if rpcResp != nil && rpcResp.Data != nil {
		ids := make([]int64, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			ids = append(ids, v.TargetId)
		}
		resp.ResourceIds = ids
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
