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

type ConcernedAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedAddLogic {
	return &ConcernedAddLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConcernedAddLogic) Add(userId int64, req *types.ConcernedAddRequest) (resp *types.ConcernedAddResponse, err error) {
	resp = new(types.ConcernedAddResponse)

	_, err = l.svcCtx.ConcernedRPC.AddConcerned(l.ctx, &concernedpb.AddConcernedRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[ConcernedAdd] rpc err: %v", err)
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
				Type:          6, // 6 是收藏
				Title:         "新收藏",
				Content:       fmt.Sprintf("用户 %s 收藏了您的作品", triggerName),
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

	return &types.ConcernedAddResponse{}, nil
}
