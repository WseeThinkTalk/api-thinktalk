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

	var items []types.SearchInfo
	if ret != nil && ret.Data != nil {
		for _, v := range ret.Data.Items {
			if v != nil {
				items = append(items, types.SearchInfo{
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
		}
		resp.Cursor = ret.Data.Cursor
		resp.IsEnd = ret.Data.IsEnd
	}

	resp.Articles = items
	return resp, nil
}
