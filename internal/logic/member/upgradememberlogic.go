package member

import (
	"context"

	"api-thinktalk/client/member/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpgradeMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpgradeMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpgradeMemberLogic {
	return &UpgradeMemberLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpgradeMemberLogic) UpgradeMember(userId int64, req *types.UpgradeMemberRequest) (*types.UpgradeMemberResponse, error) {
	_, err := l.svcCtx.MemberRPC.UpgradeMember(l.ctx, &pb.UpgradeMemberRequest{
		UserId:        userId,
		Level:         req.Level,
		DurationDays:  req.DurationDays,
		TransactionId: req.TransactionId,
	})
	if err != nil {
		l.Errorf("[UpgradeMember] rpc err: %v", err)
		return nil, err
	}
	return &types.UpgradeMemberResponse{}, nil
}
