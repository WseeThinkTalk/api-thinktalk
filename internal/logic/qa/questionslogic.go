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
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
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
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.QuestionItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
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
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
