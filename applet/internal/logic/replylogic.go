package logic

import (
	"context"

	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
	"api-thinktalk/client/reply/reply"
	"api-thinktalk/client/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyLogic {
	return &ReplyLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyLogic) CreateReply(userId int64, req *types.ReplyCreateRequest) (*types.ReplyCreateResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.CreateReply(l.ctx, &reply.CreateReplyRequest{
		BizId:         req.BizId,
		TargetId:      req.TargetId,
		ReplyUserId:   userId,
		BeReplyUserId: req.BeReplyUserId,
		ParentId:      req.ParentId,
		Content:       req.Content,
	})
	if err != nil {
		l.Errorf("[CreateReply] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyCreateResponse{ReplyId: resp.ReplyId}, nil
}

func (l *ReplyLogic) DeleteReply(userId int64, req *types.ReplyDeleteRequest) (*types.ReplyDeleteResponse, error) {
	_, err := l.svcCtx.ReplyRPC.DeleteReply(l.ctx, &reply.DeleteReplyRequest{
		ReplyId: req.ReplyId,
		UserId:  userId,
	})
	if err != nil {
		l.Errorf("[DeleteReply] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyDeleteResponse{}, nil
}

func (l *ReplyLogic) ReplyDetail(req *types.ReplyDetailRequest) (*types.ReplyDetailResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyDetail(l.ctx, &reply.ReplyDetailRequest{
		ReplyId: req.ReplyId,
	})
	if err != nil {
		l.Errorf("[ReplyDetail] rpc err: %v", err)
		return nil, err
	}
	if resp == nil || resp.Reply == nil {
		return &types.ReplyDetailResponse{}, nil
	}

	userMap := make(map[int64]*user.FindByIdResponse)
	getUserInfo := func(uid int64) (string, string) {
		if uid == 0 {
			return "匿名用户", ""
		}
		if u, ok := userMap[uid]; ok {
			return u.Username, u.Avatar
		}
		u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: uid})
		if err != nil {
			return "用户", ""
		}
		userMap[uid] = u
		return u.Username, u.Avatar
	}

	return &types.ReplyDetailResponse{Reply: convertReplyItem(resp.Reply, getUserInfo)}, nil
}

func (l *ReplyLogic) ReplyList(req *types.ReplyListRequest) (*types.ReplyListResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyList(l.ctx, &reply.ReplyListRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
	})
	if err != nil {
		l.Errorf("[ReplyList] rpc err: %v", err)
		return nil, err
	}

	userMap := make(map[int64]*user.FindByIdResponse)
	getUserInfo := func(uid int64) (string, string) {
		if uid == 0 {
			return "匿名用户", ""
		}
		if u, ok := userMap[uid]; ok {
			return u.Username, u.Avatar
		}
		u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: uid})
		if err != nil {
			return "用户", ""
		}
		userMap[uid] = u
		return u.Username, u.Avatar
	}

	items := make([]*types.ReplyItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, convertReplyItem(item, getUserInfo))
	}
	return &types.ReplyListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *ReplyLogic) ReplyCount(req *types.ReplyCountRequest) (*types.ReplyCountResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyCount(l.ctx, &reply.ReplyCountRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
	})
	if err != nil {
		l.Errorf("[ReplyCount] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyCountResponse{
		ReplyNum:     resp.ReplyNum,
		ReplyRootNum: resp.ReplyRootNum,
	}, nil
}

func convertReplyItem(pb *reply.ReplyItem, getUserInfo func(int64) (string, string)) *types.ReplyItem {
	if pb == nil {
		return nil
	}
	name, avatar := getUserInfo(pb.ReplyUserId)
	item := &types.ReplyItem{
		ReplyId:         pb.ReplyId,
		BizId:           pb.BizId,
		TargetId:        pb.TargetId,
		ReplyUserId:     pb.ReplyUserId,
		BeReplyUserId:   pb.BeReplyUserId,
		ParentId:        pb.ParentId,
		Content:         pb.Content,
		LikeNum:         pb.LikeNum,
		CreateTime:      pb.CreateTime,
		ReplyUserName:   name,
		ReplyUserAvatar: avatar,
	}
	for _, sub := range pb.SubReplies {
		item.SubReplies = append(item.SubReplies, convertReplyItem(sub, getUserInfo))
	}
	return item
}
