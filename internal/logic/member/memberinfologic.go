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

func (l *MemberInfoLogic) MemberInfo(userId int64) (*types.MemberInfoResponse, error) {
	resp, err := l.svcCtx.MemberRPC.MemberInfo(l.ctx, &pb.MemberInfoRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[MemberInfo] rpc err: %v", err)
		return nil, err
	}
	return &types.MemberInfoResponse{
		UserId:     resp.UserId,
		Level:      resp.Level,
		LevelName:  resp.LevelName,
		ExpireTime: resp.ExpireTime,
		Status:     resp.Status,
	}, nil
}
