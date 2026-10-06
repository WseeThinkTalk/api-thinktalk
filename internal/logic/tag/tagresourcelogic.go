package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagResourceLogic {
	return &TagResourceLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TagResourceLogic) TagResource(userId int64, req *types.TagResourceRequest) (resp *types.TagResourceResponse, err error) {
	resp = new(types.TagResourceResponse)

	for _, tagId := range req.TagIds {
		_, err = l.svcCtx.TagRPC.TagResource(l.ctx, &tag.TagResourceRequest{
			BizId:    req.ResourceType,
			TargetId: req.ResourceId,
			TagId:    tagId,
			UserId:   userId,
		})
		if err != nil {
			l.Errorf("[TagResource] rpc err: %v", err)
			return nil, err
		}
	}
	return resp, nil
}
