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

	// 转换评论管理列表项（支持嵌套子评论）
	var items []*types.ReplyItem
	for _, v := range rpcResp.Items {
		var subReplies []*types.ReplyItem
		for _, sub := range v.SubReplies {
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
			ReplyId:       v.ReplyId,
			BizId:         v.BizId,
			TargetId:      v.TargetId,
			ReplyUserId:   v.ReplyUserId,
			BeReplyUserId: v.BeReplyUserId,
			ParentId:      v.ParentId,
			Content:       v.Content,
			LikeNum:       v.LikeNum,
			CreateTime:    v.CreateTime,
			SubReplies:    subReplies,
		})
	}

	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
