package member

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	member "api-thinktalk/client/member/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberRightLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMemberRightLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberRightLogic {
	return &MemberRightLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MemberRightLogic) CheckRight(userId int64, req *types.MemberRightRequest) (resp *types.MemberRightResponse, err error) {
	resp = new(types.MemberRightResponse)

	rpcResp, err := l.svcCtx.MemberRPC.CheckMemberRight(l.ctx, &member.CheckMemberRightRequest{
		UserId:   userId,
		RightKey: req.RightKey,
	})
	if err != nil {
		l.Errorf("[CheckMemberRight] rpc err: %v", err)
		return nil, err
	}
	if rpcResp != nil && rpcResp.Data != nil {
		resp.HasRight = rpcResp.Data.HasRight
		resp.Level = rpcResp.Data.Level
	}
	return resp, nil
}
