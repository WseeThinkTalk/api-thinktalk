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

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FollowLogic) Follow(userId int64, req *types.FollowRequest) (resp *types.FollowResponse, err error) {
	resp = new(types.FollowResponse)

	_, err = l.svcCtx.FollowRPC.Follow(l.ctx, &pb.FollowRequest{
		UserId:         userId,
		FollowedUserId: req.FollowedUserId,
	})
	if err != nil {
		l.Errorf("[Follow] rpc err: %v", err)
		return nil, err
	}

	threading.GoSafe(func() {
		triggerName := "某用户"
		if userResp, err := l.svcCtx.UserRPC.FindById(context.Background(), &user.FindByIdRequest{UserId: userId}); err == nil {
			triggerName = userResp.Username
		}

		msg := &types.NotificationMsg{
			UserId:        req.FollowedUserId,
			Type:          3, // Follow
			Title:         "新关注",
			Content:       fmt.Sprintf("用户 %s 关注了您", triggerName),
			RefId:         userId,
			BizId:         "follow",
			TriggerUserId: userId,
		}
		data, _ := json.Marshal(msg)
		if l.svcCtx.NotificationPusher != nil {
			_ = l.svcCtx.NotificationPusher.Push(context.Background(), string(data))
		}
	})

	return resp, nil
}
