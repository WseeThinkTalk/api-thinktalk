package logic

import (
	"context"
	"fmt"

	"encoding/json"

	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
	"api-thinktalk/client/follow/pb"
	"api-thinktalk/client/user/user"

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

func (l *FollowLogic) Follow(userId int64, req *types.FollowRequest) (*types.FollowResponse, error) {
	_, err := l.svcCtx.FollowRPC.Follow(l.ctx, &pb.FollowRequest{
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

	return &types.FollowResponse{}, nil
}

func (l *FollowLogic) UnFollow(userId int64, req *types.UnfollowRequest) (*types.UnfollowResponse, error) {
	_, err := l.svcCtx.FollowRPC.UnFollow(l.ctx, &pb.UnFollowRequest{
		UserId:         userId,
		FollowedUserId: req.FollowedUserId,
	})
	if err != nil {
		l.Errorf("[UnFollow] rpc err: %v", err)
		return nil, err
	}

	threading.GoSafe(func() {
		triggerName := "某用户"
		if userResp, err := l.svcCtx.UserRPC.FindById(context.Background(), &user.FindByIdRequest{UserId: userId}); err == nil {
			triggerName = userResp.Username
		}

		msg := &types.NotificationMsg{
			UserId:        req.FollowedUserId,
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

	return &types.UnfollowResponse{}, nil
}

func (l *FollowLogic) FollowList(userId int64, req *types.FollowListRequest) (*types.FollowListResponse, error) {
	resp, err := l.svcCtx.FollowRPC.FollowList(l.ctx, &pb.FollowListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[FollowList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.FollowItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		targetName := ""
		targetAvatar := ""
		if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: item.FollowedUserId}); err == nil {
			targetName = userResp.Username
			targetAvatar = userResp.Avatar
		}
		items = append(items, &types.FollowItem{
			Id:               item.Id,
			FollowedUserId:   item.FollowedUserId,
			CreateTime:       item.CreateTime,
			FansCount:        item.FansCount,
			TargetUserName:   targetName,
			TargetUserAvatar: targetAvatar,
		})
	}
	return &types.FollowListResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}

func (l *FollowLogic) FansList(userId int64, req *types.FansListRequest) (*types.FansListResponse, error) {
	resp, err := l.svcCtx.FollowRPC.FansList(l.ctx, &pb.FansListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[FansList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.FollowItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		targetName := ""
		targetAvatar := ""
		if userResp, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: item.FansUserId}); err == nil {
			targetName = userResp.Username
			targetAvatar = userResp.Avatar
		}
		items = append(items, &types.FollowItem{
			Id:               item.UserId,
			FollowedUserId:   item.FansUserId,
			CreateTime:       item.CreateTime,
			FansCount:        item.FansCount,
			TargetUserName:   targetName,
			TargetUserAvatar: targetAvatar,
		})
	}
	return &types.FansListResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
