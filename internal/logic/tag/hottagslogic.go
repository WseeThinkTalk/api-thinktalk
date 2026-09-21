package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HotTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHotTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HotTagsLogic {
	return &HotTagsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *HotTagsLogic) HotTags(req *types.HotTagsRequest) (resp *types.HotTagsResponse, err error) {
	resp = new(types.HotTagsResponse)

	rpcResp, err := l.svcCtx.TagRPC.HotTags(l.ctx, &tag.HotTagsRequest{
		Limit: req.Limit,
	})
	if err != nil {
		l.Errorf("[HotTags] rpc err: %v", err)
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
