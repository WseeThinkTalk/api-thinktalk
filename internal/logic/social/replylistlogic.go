package social

import (
	"context"

	reply "api-thinktalk/client/reply/pb"
	user "api-thinktalk/client/user/service"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyListLogic {
	return &ReplyListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyListLogic) ReplyList(req *types.ReplyListRequest) (resp *types.ReplyListResponse, err error) {
	resp = new(types.ReplyListResponse)
	resp.Items = make([]*types.ReplyItem, 0)

	rpcResp, err := l.svcCtx.ReplyRPC.ReplyList(l.ctx, &reply.ReplyListRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
	})
	if err != nil {
		l.Errorf("[ReplyList] rpc err: %v", err)
		return nil, err
	}

	userMap := make(map[int64]*user.FindByIdResponse)
	getUserInfo := func(uid int64) (string, string) {
		if uid == 0 {
			return "匿名用户", ""
		}
		if u, ok := userMap[uid]; ok {
			return u.Username, u.Avatar
		}
		u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: uid})
		if err != nil {
			return "用户", ""
		}
		userMap[uid] = u
		return u.Username, u.Avatar
	}

	items := make([]*types.ReplyItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		items = append(items, convertReplyItem(item, getUserInfo))
	}
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
