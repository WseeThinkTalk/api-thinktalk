package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagUpdateLogic {
	return &TagUpdateLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TagUpdateLogic) TagUpdate(userId int64, req *types.TagUpdateRequest) (resp *types.TagUpdateResponse, err error) {
	resp = new(types.TagUpdateResponse)

	_, err = l.svcCtx.TagRPC.UpdateTag(l.ctx, &tag.UpdateTagRequest{
		TagId:   req.TagId,
		TagName: req.TagName,
		TagDesc: req.TagDesc,
	})
	if err != nil {
		l.Errorf("[UpdateTag] rpc err: %v", err)
		return nil, err
	}
	return resp, nil
}
