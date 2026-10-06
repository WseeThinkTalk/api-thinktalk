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
	resp = new(types.AdminPendingListResponse)

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

	var items []types.SearchInfo
	if ret != nil && ret.Data != nil {
		for _, v := range ret.Data.Items {
			if v != nil {
				items = append(items, types.SearchInfo{
					ArticleId:   v.ArticleId,
					Title:       v.Title,
					Description: v.Description,
					Cover:       v.Cover,
					AuthorId:    v.AuthorId,
					AuthorName:  v.AuthorName,
					LikeNum:     v.LikeNum,
					CommentNum:  v.CommentNum,
					PublishTime: v.PublishTime,
					AuthorAvatar: v.AuthorAvatar,
				})
			}
		}
		resp.Cursor = ret.Data.Cursor
		resp.IsEnd = ret.Data.IsEnd
	}

	resp.Articles = items
	return resp, nil
}
