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
		BizId:    req.BizId,
		TargetId: req.TargetId,
	})
	if err != nil {
		l.Errorf("[TagsByResource] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.TagItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		items = append(items, &types.TagItem{
			TagId:         item.TagId,
			TagName:       item.TagName,
			TagDesc:       item.TagDesc,
			ResourceCount: item.ResourceCount,
			CreateTime:    item.CreateTime,
		})
	}
	resp.Items = items
	return resp, nil
}
