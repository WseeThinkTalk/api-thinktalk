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

func (l *TagsByResourceLogic) TagsByResource(req *types.TagsByResourceRequest) (*types.TagsByResourceResponse, error) {
	resp, err := l.svcCtx.TagRPC.TagsByResource(l.ctx, &tag.TagsByResourceRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
	})
	if err != nil {
		l.Errorf("[TagsByResource] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.TagItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.TagItem{
			TagId:         item.TagId,
			TagName:       item.TagName,
			TagDesc:       item.TagDesc,
			ResourceCount: item.ResourceCount,
			CreateTime:    item.CreateTime,
		})
	}
	return &types.TagsByResourceResponse{Items: items}, nil
}
