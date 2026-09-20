package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionsLogic {
	return &QuestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuestionsLogic) Questions(userId int64, req *types.QuestionListRequest) (*types.QuestionListResponse, error) {
	resp, err := l.svcCtx.QaRPC.Questions(l.ctx, &qa.QuestionsRequest{
		UserId:   req.AuthorId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
	})
	if err != nil {
		l.Errorf("[Questions] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.QuestionItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.QuestionItem{
			Id:         item.Id,
			Title:      item.Title,
			Content:    item.Content,
			AuthorId:   item.AuthorId,
			AnswerNum:  item.AnswerNum,
			ViewNum:    item.ViewNum,
			TagIds:     item.TagIds,
			CreateTime: item.CreateTime,
		})
	}
	return &types.QuestionListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}
