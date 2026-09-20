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

func (l *TagListLogic) TagList(req *types.TagListRequest) (*types.TagListResponse, error) {
	resp, err := l.svcCtx.TagRPC.TagList(l.ctx, &tag.TagListRequest{
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[TagList] rpc err: %v", err)
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
	return &types.TagListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}
