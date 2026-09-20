// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	article "api-thinktalk/client/article/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminPendingListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminPendingListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminPendingListLogic {
	return &AdminPendingListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminPendingListLogic) AdminPendingList(req *types.AdminPendingListRequest) (resp *types.AdminPendingListResponse, err error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	ret, err := l.svcCtx.ArticleRPC.AdminPendingList(l.ctx, &article.AdminPendingListRequest{
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("AdminPendingList error: %v", err)
		return nil, err
	}

	items := make([]types.SearchInfo, 0, len(ret.Items))
	for _, item := range ret.Items {
		items = append(items, types.SearchInfo{
			ArticleId:   item.ArticleId,
			Title:       item.Title,
			Description: item.Description,
			Cover:       item.Cover,
			AuthorId:    item.AuthorId,
			AuthorName:  item.AuthorName,
			LikeNum:     item.LikeNum,
			CommentNum:  item.CommentNum,
			PublishTime: item.PublishTime,
			AuthorAvatar: item.AuthorAvatar,
		})
	}

	return &types.AdminPendingListResponse{
		Articles: items,
		Cursor:   ret.Cursor,
		IsEnd:    ret.IsEnd,
	}, nil
}
