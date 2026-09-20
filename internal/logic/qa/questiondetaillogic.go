package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuestionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionDetailLogic {
	return &QuestionDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuestionDetailLogic) QuestionDetail(req *types.QuestionDetailRequest) (*types.QuestionDetailResponse, error) {
	resp, err := l.svcCtx.QaRPC.QuestionDetail(l.ctx, &qa.QuestionDetailRequest{
		QuestionId: req.QuestionId,
	})
	if err != nil {
		l.Errorf("[QuestionDetail] rpc err: %v", err)
		return nil, err
	}
	if resp == nil || resp.Question == nil {
		return &types.QuestionDetailResponse{}, nil
	}
	q := resp.Question
	return &types.QuestionDetailResponse{
		Question: &types.QuestionItem{
			Id:         q.Id,
			Title:      q.Title,
			Content:    q.Content,
			AuthorId:   q.AuthorId,
			AnswerNum:  q.AnswerNum,
			ViewNum:    q.ViewNum,
			TagIds:     q.TagIds,
			CreateTime: q.CreateTime,
		},
	}, nil
}
