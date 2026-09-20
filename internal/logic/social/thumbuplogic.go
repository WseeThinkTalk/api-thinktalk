package social

import (
	"context"
	"fmt"
	"encoding/json"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	"api-thinktalk/client/article/pb"
	like "api-thinktalk/client/like/service"
	replypb "api-thinktalk/client/reply/pb"
	user "api-thinktalk/client/user/service"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type ThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThumbupLogic) Thumbup(req *types.ThumbupRequest) (*types.ThumbupResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.LikeRPC.Thumbup(l.ctx, &like.ThumbupRequest{
		BizId:    req.BizId,
		ObjId:    req.ObjId,
		UserId:   userId,
		LikeType: req.LikeType,
	})
	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		var targetUserId int64
		if req.BizId == "article" {
			artResp, err := l.svcCtx.ArticleRPC.ArticleDetail(context.Background(), &pb.ArticleDetailRequest{ArticleId: req.ObjId})
			if err == nil && artResp.Article != nil {
				targetUserId = artResp.Article.AuthorId
			}
		} else if req.BizId == "reply" {
			repResp, err := l.svcCtx.ReplyRPC.ReplyDetail(context.Background(), &replypb.ReplyDetailRequest{ReplyId: req.ObjId})
			if err == nil && repResp.Reply != nil {
				targetUserId = repResp.Reply.ReplyUserId
			}
		}

		if targetUserId > 0 {
			triggerName := "某用户"
			if userResp, err := l.svcCtx.UserRPC.FindById(context.Background(), &user.FindByIdRequest{UserId: userId}); err == nil {
				triggerName = userResp.Username
			}

			msgType := int32(1)
			title := "新点赞"
			content := fmt.Sprintf("用户 %s 点赞了您", triggerName)
			if req.LikeType == 0 { // 假设 0 是取消点赞
				msgType = 5
				title = "取消点赞"
				content = fmt.Sprintf("用户 %s 取消了点赞", triggerName)
			}
			msg := &types.NotificationMsg{
				UserId:        targetUserId,
				Type:          msgType,
				Title:         title,
				Content:       content,
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

	return &types.ThumbupResponse{
		BizId:      resp.BizId,
		ObjId:      resp.ObjId,
		LikeNum:    resp.LikeNum,
		DislikeNum: resp.DislikeNum,
	}, nil
}
