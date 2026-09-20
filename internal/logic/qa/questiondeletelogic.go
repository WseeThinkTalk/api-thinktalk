package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuestionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionDeleteLogic {
	return &QuestionDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuestionDeleteLogic) QuestionDelete(userId int64, req *types.QuestionDeleteRequest) error {
	_, err := l.svcCtx.QaRPC.QuestionDelete(l.ctx, &qa.QuestionDeleteRequest{
		UserId:     userId,
		QuestionId: req.QuestionId,
	})
	if err != nil {
		l.Errorf("[QuestionDelete] rpc err: %v", err)
		return err
	}
	return nil
}
