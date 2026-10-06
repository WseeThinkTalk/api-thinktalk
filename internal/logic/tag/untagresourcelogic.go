package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UntagResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUntagResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UntagResourceLogic {
	return &UntagResourceLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UntagResourceLogic) UntagResource(userId int64, req *types.UntagResourceRequest) (resp *types.UntagResourceResponse, err error) {
	resp = new(types.UntagResourceResponse)

	for _, tagId := range req.TagIds {
		_, err = l.svcCtx.TagRPC.UntagResource(l.ctx, &tag.UntagResourceRequest{
			BizId:    req.ResourceType,
			TargetId: req.ResourceId,
			TagId:    tagId,
			UserId:   userId,
		})
		if err != nil {
			l.Errorf("[UntagResource] rpc err: %v", err)
			return nil, err
		}
	}
	return resp, nil
}
