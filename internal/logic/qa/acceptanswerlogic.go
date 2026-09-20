package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AcceptAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAcceptAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptAnswerLogic {
	return &AcceptAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AcceptAnswerLogic) AcceptAnswer(userId int64, req *types.AcceptAnswerRequest) error {
	_, err := l.svcCtx.QaRPC.AcceptAnswer(l.ctx, &qa.AcceptAnswerRequest{
		QuestionId: req.QuestionId,
		AnswerId:   req.AnswerId,
		UserId:     userId,
	})
	if err != nil {
		l.Errorf("[AcceptAnswer] rpc err: %v", err)
		return err
	}
	return nil
}
