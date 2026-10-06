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
	// 转换问答搜索结果项
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.SearchQuestionItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			items = append(items, &types.SearchQuestionItem{
				Id:         v.Id,
				Title:      v.Title,
				Content:    v.Content,
				AuthorId:   v.AuthorId,
				AnswerNum:  v.AnswerNum,
				TagIds:     v.TagIds,
				CreateTime: v.CreateTime,
			})
		}
		resp.Items = items
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
