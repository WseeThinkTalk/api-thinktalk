package social

import (
	"context"

	reply "api-thinktalk/client/reply/pb"
	user "api-thinktalk/client/user/service"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyDetailLogic {
	return &ReplyDetailLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyDetailLogic) ReplyDetail(req *types.ReplyDetailRequest) (resp *types.ReplyDetailResponse, err error) {
	resp = new(types.ReplyDetailResponse)

	rpcResp, err := l.svcCtx.ReplyRPC.ReplyDetail(l.ctx, &reply.ReplyDetailRequest{
		ReplyId: req.ReplyId,
	})
	if err != nil {
		l.Errorf("[ReplyDetail] rpc err: %v", err)
		return nil, err
	}
	if rpcResp == nil || rpcResp.Reply == nil {
		return resp, nil
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

	resp.Reply = convertReplyItem(rpcResp.Reply, getUserInfo)
	return resp, nil
}
