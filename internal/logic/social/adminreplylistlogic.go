package social

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	reply "api-thinktalk/client/reply/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminReplyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminReplyListLogic {
	return &AdminReplyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminReplyListLogic) AdminReplyList(req *types.AdminReplyListRequest) (resp *types.AdminReplyListResponse, err error) {
	resp = new(types.AdminReplyListResponse)

	rpcResp, err := l.svcCtx.ReplyRPC.AdminReplyList(l.ctx, &reply.AdminReplyListRequest{
		Keyword:  "",
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	getUserInfo := func(uid int64) (string, string) {
		return "", ""
	}

	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.ReplyItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			items = append(items, convertReplyItem(v, getUserInfo))
		}
		resp.Items = items
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
