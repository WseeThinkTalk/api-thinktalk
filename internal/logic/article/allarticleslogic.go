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

	if ret != nil && ret.Data != nil {
		resp.Articles = convertSearchItems(ret.Data.Items)
		resp.Cursor = ret.Data.Cursor
		resp.IsEnd = ret.Data.IsEnd
	}
	return resp, nil
}

func convertSearchItems(items []*article.SearchItem) []types.SearchInfo {
	if len(items) == 0 {
		return nil
	}
	res := make([]types.SearchInfo, 0, len(items))
	for _, v := range items {
		if v == nil {
			continue
		}
		res = append(res, types.SearchInfo{
			ArticleId:    v.ArticleId,
			Title:        v.Title,
			Description:  v.Description,
			Cover:        v.Cover,
			AuthorId:     v.AuthorId,
			AuthorName:   v.AuthorName,
			LikeNum:      v.LikeNum,
			CommentNum:   v.CommentNum,
			PublishTime:  v.PublishTime,
			AuthorAvatar: v.AuthorAvatar,
		})
	}
	return res
}
