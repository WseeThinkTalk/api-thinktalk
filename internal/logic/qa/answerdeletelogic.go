package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerDeleteLogic {
	return &AnswerDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnswerDeleteLogic) AnswerDelete(userId int64, req *types.AnswerDeleteRequest) error {
	_, err := l.svcCtx.QaRPC.AnswerDelete(l.ctx, &qa.AnswerDeleteRequest{
		UserId:   userId,
		AnswerId: req.AnswerId,
	})
	if err != nil {
		l.Errorf("[AnswerDelete] rpc err: %v", err)
		return err
	}
	return nil
}
