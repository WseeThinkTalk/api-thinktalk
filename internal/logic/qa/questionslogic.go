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

func (l *QuestionsLogic) Questions(userId int64, req *types.QuestionListRequest) (resp *types.QuestionListResponse, err error) {
	resp = new(types.QuestionListResponse)

	rpcResp, err := l.svcCtx.QaRPC.Questions(l.ctx, &qa.QuestionsRequest{
		UserId:   req.AuthorId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
	})
	if err != nil {
		l.Errorf("[Questions] rpc err: %v", err)
		return nil, err
	}
	// 转换问答列表项
	items := make([]*types.QuestionItem, 0, len(rpcResp.Items))
	for _, v := range rpcResp.Items {
		items = append(items, &types.QuestionItem{
			Id:         v.Id,
			Title:      v.Title,
			Content:    v.Content,
			AuthorId:   v.AuthorId,
			AnswerNum:  v.AnswerNum,
			ViewNum:    v.ViewNum,
			TagIds:     v.TagIds,
			CreateTime: v.CreateTime,
		})
	}
	resp.Items = items
	resp.Cursor = rpcResp.Cursor
	resp.IsEnd = rpcResp.IsEnd
	return resp, nil
}
