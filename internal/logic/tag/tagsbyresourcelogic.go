package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagsByResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagsByResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagsByResourceLogic {
	return &TagsByResourceLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TagsByResourceLogic) TagsByResource(req *types.TagsByResourceRequest) (resp *types.TagsByResourceResponse, err error) {
	resp = new(types.TagsByResourceResponse)

	rpcResp, err := l.svcCtx.TagRPC.TagsByResource(l.ctx, &tag.TagsByResourceRequest{
		BizId:    req.ResourceType,
		TargetId: req.ResourceId,
	})
	if err != nil {
		l.Errorf("[TagsByResource] rpc err: %v", err)
		return nil, err
	}
	// 转换资源标签列表数据项
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.TagItem, 0, len(rpcResp.Data))
		for _, v := range rpcResp.Data {
			items = append(items, &types.TagItem{
				Id:          v.TagId,
				Name:        v.TagName,
				Description: v.TagDesc,
				UseCount:    v.ResourceCount,
			})
		}
		resp.Tags = items
	}
	return resp, nil
}
