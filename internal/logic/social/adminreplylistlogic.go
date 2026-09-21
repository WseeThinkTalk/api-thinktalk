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
		Keyword:  req.Keyword,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var items []*types.ReplyItem
	for _, item := range rpcResp.Items {
		var subReplies []*types.ReplyItem
		for _, sub := range item.SubReplies {
			subReplies = append(subReplies, &types.ReplyItem{
				ReplyId:       sub.ReplyId,
				BizId:         sub.BizId,
				TargetId:      sub.TargetId,
				ReplyUserId:   sub.ReplyUserId,
				BeReplyUserId: sub.BeReplyUserId,
				ParentId:      sub.ParentId,
				Content:       sub.Content,
				LikeNum:       sub.LikeNum,
				CreateTime:    sub.CreateTime,
			})
		}

		items = append(items, &types.ReplyItem{
			ReplyId:       item.ReplyId,
			BizId:         item.BizId,
			TargetId:      item.TargetId,
			ReplyUserId:   item.ReplyUserId,
			BeReplyUserId: item.BeReplyUserId,
			ParentId:      item.ParentId,
			Content:       item.Content,
			LikeNum:       item.LikeNum,
			CreateTime:    item.CreateTime,
			SubReplies:    subReplies,
		})
	}

	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
