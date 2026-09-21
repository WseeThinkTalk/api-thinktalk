package member

import (
	"context"

	"api-thinktalk/client/member/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

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

	rpcResp, err := l.svcCtx.MemberRPC.MemberInfo(l.ctx, &pb.MemberInfoRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[MemberInfo] rpc err: %v", err)
		return nil, err
	}
	resp.UserId = rpcResp.UserId
	resp.Level = rpcResp.Level
	resp.LevelName = rpcResp.LevelName
	resp.ExpireTime = rpcResp.ExpireTime
	resp.Status = rpcResp.Status
	return resp, nil
}
