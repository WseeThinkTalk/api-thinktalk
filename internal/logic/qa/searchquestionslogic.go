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

func (l *SearchQuestionsLogic) SearchQuestions(req *types.SearchQuestionsRequest) (resp *types.SearchQuestionsResponse, err error) {
	resp = new(types.SearchQuestionsResponse)

	rpcResp, err := l.svcCtx.QaRPC.SearchQuestions(l.ctx, &qa.SearchQuestionsRequest{
		Keyword:  req.Keyword,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[SearchQuestions] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.SearchQuestionItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
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
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
