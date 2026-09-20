package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchQuestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchQuestionsLogic {
	return &SearchQuestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchQuestionsLogic) SearchQuestions(req *types.SearchQuestionsRequest) (*types.SearchQuestionsResponse, error) {
	resp, err := l.svcCtx.QaRPC.SearchQuestions(l.ctx, &qa.SearchQuestionsRequest{
		Keyword:  req.Keyword,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[SearchQuestions] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.SearchQuestionItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.SearchQuestionItem{
			Id:         item.Id,
			Title:      item.Title,
			Content:    item.Content,
			AuthorId:   item.AuthorId,
			AnswerNum:  item.AnswerNum,
			TagIds:     item.TagIds,
			CreateTime: item.CreateTime,
		})
	}
	return &types.SearchQuestionsResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}
