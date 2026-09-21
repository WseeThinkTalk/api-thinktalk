package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerQuestionLogic {
	return &AnswerQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnswerQuestionLogic) AnswerQuestion(userId int64, req *types.AnswerQuestionRequest) (resp *types.AnswerQuestionResponse, err error) {
	resp = new(types.AnswerQuestionResponse)

	rpcResp, err := l.svcCtx.QaRPC.AnswerQuestion(l.ctx, &qa.AnswerQuestionRequest{
		QuestionId: req.QuestionId,
		UserId:     userId,
		Content:    req.Content,
	})
	if err != nil {
		l.Errorf("[AnswerQuestion] rpc err: %v", err)
		return nil, err
	}
	resp.AnswerId = rpcResp.AnswerId
	return resp, nil
}
