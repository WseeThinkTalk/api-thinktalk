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

func (l *AnswerListLogic) AnswerList(req *types.AnswerListRequest) (resp *types.AnswerListResponse, err error) {
	resp = new(types.AnswerListResponse)

	rpcResp, err := l.svcCtx.QaRPC.AnswerList(l.ctx, &qa.AnswerListRequest{
		QuestionId: req.QuestionId,
		Cursor:     req.Cursor,
		PageSize:   req.PageSize,
	})
	if err != nil {
		l.Errorf("[AnswerList] rpc err: %v", err)
		return nil, err
	}
	// 转换回答列表项
	if rpcResp != nil && rpcResp.Data != nil {
		items := make([]*types.AnswerItem, 0, len(rpcResp.Data.Items))
		for _, v := range rpcResp.Data.Items {
			items = append(items, &types.AnswerItem{
				Id:         v.Id,
				QuestionId: v.QuestionId,
				AuthorId:   v.AuthorId,
				Content:    v.Content,
				IsAccepted: v.IsAccepted,
				LikeNum:    v.LikeNum,
				ReplyNum:   v.ReplyNum,
				CreateTime: v.CreateTime,
			})
		}
		resp.Items = items
		resp.Cursor = rpcResp.Data.Cursor
		resp.IsEnd = rpcResp.Data.IsEnd
	}
	return resp, nil
}
