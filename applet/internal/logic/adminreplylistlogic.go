package logic

import (
	"context"

	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
	"api-thinktalk/client/reply/reply"

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

	return &types.AdminReplyListResponse{
		Items:  items,
		Cursor: rpcResp.Cursor,
		IsEnd:  rpcResp.IsEnd,
	}, nil
}
