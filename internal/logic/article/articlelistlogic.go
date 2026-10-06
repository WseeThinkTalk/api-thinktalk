// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"
	"encoding/json"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	article "api-thinktalk/client/article/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleListLogic {
	return &ArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleListLogic) ArticleList(req *types.ArticleListRequest) (resp *types.ArticleListResponse, err error) {
	resp = new(types.ArticleListResponse)

	userId := req.AuthorId
	if userId <= 0 {
		if val := l.ctx.Value("userId"); val != nil {
			if num, ok := val.(json.Number); ok {
				userId, _ = num.Int64()
			}
		}
	}

	articles, err := l.svcCtx.ArticleRPC.Articles(l.ctx, &article.ArticlesRequest{
		UserId:    userId,
		Cursor:    req.Cursor,
		PageSize:  req.PageSize,
		SortType:  req.SortType,
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("get articles req: %v err: %v", req, err)
		return nil, err
	}
	if articles == nil || articles.Data == nil || len(articles.Data.Articles) == 0 {
		resp.Articles = make([]types.ArticleInfo, 0)
		return resp, nil
	}

	// 转换用户文章列表数据项
	infos := make([]types.ArticleInfo, 0, len(articles.Data.Articles))
	for _, v := range articles.Data.Articles {
		if v != nil {
			infos = append(infos, types.ArticleInfo{
				ArticleId:   v.Id,
				Cover:       v.Cover,
				Description: v.Description,
				Title:       v.Title,
				Status:      v.Status,
			})
		}
	}

	resp.Articles = infos
	return resp, nil
}
