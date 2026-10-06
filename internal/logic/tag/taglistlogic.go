package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagListLogic {
	return &TagListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TagListLogic) TagList(req *types.TagListRequest) (resp *types.TagListResponse, err error) {
	resp = new(types.TagListResponse)

	rpcResp, err := l.svcCtx.TagRPC.TagList(l.ctx, &tag.TagListRequest{
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[TagList] rpc err: %v", err)
		return nil, err
	}
	// 转换标签列表数据项
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.TagItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			items = append(items, &types.TagItem{
				Id:          v.TagId,
				Name:        v.TagName,
				Description: v.TagDesc,
				UseCount:    v.ResourceCount,
			})
		}
		resp.Items = items
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
