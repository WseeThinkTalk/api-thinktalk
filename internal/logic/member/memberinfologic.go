package member

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	member "api-thinktalk/client/member/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberInfoLogic {
	return &MemberInfoLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MemberInfoLogic) MemberInfo(userId int64) (resp *types.MemberInfoResponse, err error) {
	resp = new(types.MemberInfoResponse)

	rpcResp, err := l.svcCtx.MemberRPC.MemberInfo(l.ctx, &member.MemberInfoRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[MemberInfo] rpc err: %v", err)
		return nil, err
	}
	if rpcResp != nil && rpcResp.Data != nil {
		resp.UserId = rpcResp.Data.UserId
		resp.MemberLevel = rpcResp.Data.Level
		resp.ExpireAt = rpcResp.Data.ExpireTime
		resp.IsActive = rpcResp.Data.Status == 1
	}
	return resp, nil
}
