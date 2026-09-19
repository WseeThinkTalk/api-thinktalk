package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"api-thinktalk/applet/internal/svc"
	"api-thinktalk/applet/internal/types"
	"api-thinktalk/client/article/pb"
	concernedpb "api-thinktalk/client/concerned/pb"
	"api-thinktalk/client/user/user"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type ConcernedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedLogic {
	return &ConcernedLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConcernedLogic) Add(userId int64, req *types.ConcernedAddRequest) (*types.ConcernedAddResponse, error) {
	_, err := l.svcCtx.ConcernedRPC.AddConcerned(l.ctx, &concernedpb.AddConcernedRequest{
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

func (l *ConcernedLogic) Cancel(userId int64, req *types.ConcernedCancelRequest) (*types.ConcernedCancelResponse, error) {
	_, err := l.svcCtx.ConcernedRPC.CancelConcerned(l.ctx, &concernedpb.CancelConcernedRequest{
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

	return &types.ConcernedCancelResponse{}, nil
}

func (l *ConcernedLogic) Check(userId int64, req *types.ConcernedCheckRequest) (*types.ConcernedCheckResponse, error) {
	resp, err := l.svcCtx.ConcernedRPC.IsConcerned(l.ctx, &concernedpb.IsConcernedRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[ConcernedCheck] rpc err: %v", err)
		return nil, err
	}
	return &types.ConcernedCheckResponse{IsConcerned: resp.IsConcerned}, nil
}

func (l *ConcernedLogic) ConcernedCount(req *types.ConcernedCountRequest) (*types.ConcernedCountResponse, error) {
	objId := req.ObjId
	if objId == 0 {
		userIdVal := l.ctx.Value("userId")
		if userIdVal != nil {
			if num, ok := userIdVal.(json.Number); ok {
				uid, _ := num.Int64()
				objId = -uid
			}
		}
	}

	resp, err := l.svcCtx.ConcernedRPC.ConcernedCount(l.ctx, &concernedpb.ConcernedCountRequest{
		BizId: req.BizId,
		ObjId: objId,
	})
	if err != nil {
		l.Errorf("[ConcernedCount] rpc err: %v", err)
		return nil, err
	}
	return &types.ConcernedCountResponse{ConcernedNum: resp.ConcernedNum}, nil
}

func (l *ConcernedLogic) List(userId int64, req *types.ConcernedListRequest) (*types.ConcernedListResponse, error) {
	resp, err := l.svcCtx.ConcernedRPC.ConcernedList(l.ctx, &concernedpb.ConcernedListRequest{
		UserId:   userId,
		BizId:    req.BizId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[ConcernedList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.ConcernedItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.ConcernedItem{
			Id:         item.Id,
			BizId:      item.BizId,
			ObjId:      item.ObjId,
			CreateTime: item.CreateTime,
		})
	}
	return &types.ConcernedListResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
