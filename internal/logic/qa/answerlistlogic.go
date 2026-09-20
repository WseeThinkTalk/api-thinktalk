package qa

import (
	"context"

	qa "api-thinktalk/client/qa/pb"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerListLogic {
	return &AnswerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnswerListLogic) AnswerList(req *types.AnswerListRequest) (*types.AnswerListResponse, error) {
	resp, err := l.svcCtx.QaRPC.AnswerList(l.ctx, &qa.AnswerListRequest{
		QuestionId: req.QuestionId,
		Cursor:     req.Cursor,
		PageSize:   req.PageSize,
	})
	if err != nil {
		l.Errorf("[AnswerList] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.AnswerItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.AnswerItem{
			Id:         item.Id,
			QuestionId: item.QuestionId,
			AuthorId:   item.AuthorId,
			Content:    item.Content,
			IsAccepted: item.IsAccepted,
			LikeNum:    item.LikeNum,
			ReplyNum:   item.ReplyNum,
			CreateTime: item.CreateTime,
		})
	}
	return &types.AnswerListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}
