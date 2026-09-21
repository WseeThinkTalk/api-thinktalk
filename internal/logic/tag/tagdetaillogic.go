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

func (l *TagDetailLogic) TagDetail(req *types.TagDetailRequest) (resp *types.TagDetailResponse, err error) {
	resp = new(types.TagDetailResponse)

	rpcResp, err := l.svcCtx.TagRPC.TagDetail(l.ctx, &tag.TagDetailRequest{
		TagId: req.TagId,
	})
	if err != nil {
		l.Errorf("[TagDetail] rpc err: %v", err)
		return nil, err
	}
	resp.TagId = rpcResp.TagId
	resp.TagName = rpcResp.TagName
	resp.TagDesc = rpcResp.TagDesc
	resp.ResourceCount = rpcResp.ResourceCount
	resp.CreateTime = rpcResp.CreateTime
	return resp, nil
}
