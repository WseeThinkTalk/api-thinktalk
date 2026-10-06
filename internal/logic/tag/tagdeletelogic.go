package tag

import (
	"context"

	tag "api-thinktalk/client/tag/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagDeleteLogic {
	return &TagDeleteLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TagDeleteLogic) TagDelete(userId int64, req *types.TagDeleteRequest) (resp *types.TagDeleteResponse, err error) {
	resp = new(types.TagDeleteResponse)

	_, err = l.svcCtx.TagRPC.DeleteTag(l.ctx, &tag.DeleteTagRequest{
		TagId: req.Id,
	})
	if err != nil {
		l.Errorf("[DeleteTag] rpc err: %v", err)
		return nil, err
	}
	return resp, nil
}
