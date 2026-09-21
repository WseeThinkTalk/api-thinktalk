package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagCreateLogic {
	return &TagCreateLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TagCreateLogic) TagCreate(userId int64, req *types.TagCreateRequest) (resp *types.TagCreateResponse, err error) {
	resp = new(types.TagCreateResponse)

	rpcResp, err := l.svcCtx.TagRPC.CreateTag(l.ctx, &tag.CreateTagRequest{
		TagName: req.TagName,
		TagDesc: req.TagDesc,
	})
	if err != nil {
		l.Errorf("[CreateTag] rpc err: %v", err)
		return nil, err
	}
	resp.TagId = rpcResp.TagId
	return resp, nil
}
