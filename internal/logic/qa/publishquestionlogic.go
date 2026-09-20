package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishQuestionLogic {
	return &PublishQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishQuestionLogic) PublishQuestion(userId int64, req *types.PublishQuestionRequest) (*types.PublishQuestionResponse, error) {
	resp, err := l.svcCtx.QaRPC.PublishQuestion(l.ctx, &qa.PublishQuestionRequest{
		UserId:  userId,
		Title:   req.Title,
		Content: req.Content,
		TagIds:  req.TagIds,
	})
	if err != nil {
		l.Errorf("[PublishQuestion] rpc err: %v", err)
		return nil, err
	}
	return &types.PublishQuestionResponse{QuestionId: resp.QuestionId}, nil
}
