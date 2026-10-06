package social

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/client/follow/pb"
	user "api-thinktalk/client/user/service"
	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type UnFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnFollowLogic) UnFollow(userId int64, req *types.UnFollowRequest) (resp *types.UnFollowResponse, err error) {
	resp = new(types.UnFollowResponse)

	_, err = l.svcCtx.FollowRPC.UnFollow(l.ctx, &pb.UnFollowRequest{
		UserId:         userId,
		FollowedUserId: req.FollowId,
	})
	if err != nil {
		l.Errorf("[UnFollow] rpc err: %v", err)
		return nil, err
	}

	threading.GoSafe(func() {
		triggerName := "某用户"
		if userResp, err := l.svcCtx.UserRPC.FindById(context.Background(), &user.FindByIdRequest{UserId: userId}); err == nil && userResp != nil && userResp.Data != nil {
			triggerName = userResp.Data.Username
		}

		msg := &types.NotificationMsg{
			UserId:        req.FollowId,
			Type:          4, // UnFollow
			Title:         "取消关注",
			Content:       fmt.Sprintf("用户 %s 取消了关注", triggerName),
			RefId:         userId,
			BizId:         "unfollow",
			TriggerUserId: userId,
		}
		data, _ := json.Marshal(msg)
		if l.svcCtx.NotificationPusher != nil {
			_ = l.svcCtx.NotificationPusher.Push(context.Background(), string(data))
		}
	})

	return resp, nil
}
