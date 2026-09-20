package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagDetailLogic {
	return &TagDetailLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TagDetailLogic) TagDetail(req *types.TagDetailRequest) (*types.TagDetailResponse, error) {
	resp, err := l.svcCtx.TagRPC.TagDetail(l.ctx, &tag.TagDetailRequest{
		TagId: req.TagId,
	})
	if err != nil {
		l.Errorf("[TagDetail] rpc err: %v", err)
		return nil, err
	}
	return &types.TagDetailResponse{
		TagId:         resp.TagId,
		TagName:       resp.TagName,
		TagDesc:       resp.TagDesc,
		ResourceCount: resp.ResourceCount,
		CreateTime:    resp.CreateTime,
	}, nil
}
