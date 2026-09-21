package article

import (
	"context"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	article "api-thinktalk/client/article/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	resp = new(types.SearchResponse)

	ret, err := l.svcCtx.ArticleRPC.SearchArticles(l.ctx, &article.SearchRequest{
		Keyword:  req.Keyword,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("SearchArticles req: %v error: %v", req, err)
		return nil, err
	}

	// 转换文章搜索结果数据项
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
