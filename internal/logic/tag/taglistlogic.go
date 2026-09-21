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
	items := make([]*types.TagItem, 0, len(rpcResp.Items))
	for _, v := range rpcResp.Items {
		items = append(items, &types.TagItem{
			TagId:         v.TagId,
			TagName:       v.TagName,
			TagDesc:       v.TagDesc,
			ResourceCount: v.ResourceCount,
			CreateTime:    v.CreateTime,
		})
	}
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
