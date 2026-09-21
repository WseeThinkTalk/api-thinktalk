package social

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/client/article/pb"
	concernedpb "api-thinktalk/client/concerned/pb"
	user "api-thinktalk/client/user/service"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type ConcernedCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedCancelLogic {
	return &ConcernedCancelLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConcernedCancelLogic) Cancel(userId int64, req *types.ConcernedCancelRequest) (resp *types.ConcernedCancelResponse, err error) {
	resp = new(types.ConcernedCancelResponse)

	_, err = l.svcCtx.ConcernedRPC.CancelConcerned(l.ctx, &concernedpb.CancelConcernedRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[ConcernedCancel] rpc err: %v", err)
		return nil, err
	}

	threading.GoSafe(func() {
		var targetUserId int64
		if req.BizId == "article" {
			artResp, err := l.svcCtx.ArticleRPC.ArticleDetail(context.Background(), &pb.ArticleDetailRequest{ArticleId: req.ObjId})
			if err == nil && artResp.Article != nil {
				targetUserId = artResp.Article.AuthorId
			}
		}

		if targetUserId > 0 {
			triggerName := "某用户"
			if userResp, err := l.svcCtx.UserRPC.FindById(context.Background(), &user.FindByIdRequest{UserId: userId}); err == nil {
				triggerName = userResp.Username
			}

			msg := &types.NotificationMsg{
				UserId:        targetUserId,
				Type:          7, // 7 是取消收藏
				Title:         "取消收藏",
				Content:       fmt.Sprintf("用户 %s 取消收藏了您的作品", triggerName),
				RefId:         req.ObjId,
				BizId:         req.BizId,
				TriggerUserId: userId,
			}
			data, _ := json.Marshal(msg)
			if l.svcCtx.NotificationPusher != nil {
				_ = l.svcCtx.NotificationPusher.Push(context.Background(), string(data))
			}
		}
	})

	return resp, nil
}
