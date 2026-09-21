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

type AllArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllArticlesLogic {
	return &AllArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllArticlesLogic) AllArticles(req *types.AllArticlesRequest) (resp *types.AllArticlesResponse, err error) {
	resp = new(types.AllArticlesResponse)

	if req.PageSize == 0 {
		req.PageSize = 20
	}

	ret, err := l.svcCtx.ArticleRPC.SearchArticles(l.ctx, &article.SearchRequest{
		Keyword:  "", // Empty keyword gets all articles sorted by time
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("AllArticles SearchArticles error: %v", err)
		return nil, err
	}

	// 转换全站文章数据项
	items := make([]types.SearchInfo, 0, len(ret.Items))
	for _, v := range ret.Items {
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

	resp.Articles = items
	resp.Cursor = ret.Cursor
	resp.IsEnd = ret.IsEnd
	return resp, nil
}
