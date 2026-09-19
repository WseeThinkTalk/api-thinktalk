// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"api-thinktalk/client/article/article"
	"api-thinktalk/client/user/user"
	"context"
	"time"

	"api-thinktalk/article/internal/svc"
	"api-thinktalk/article/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDetailLogic) ArticleDetail(req *types.ArticleDetailRequest) (resp *types.ArticleDetailResponse, err error) {
	articleInfo, err := l.svcCtx.ArticleRPC.ArticleDetail(l.ctx, &article.ArticleDetailRequest{
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("get article detail id: %d err: %v", req.ArticleId, err)
		return nil, err
	}
	if articleInfo == nil || articleInfo.Article == nil {
		return nil, nil
	}
	userInfo, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: articleInfo.Article.AuthorId,
	})
	if err != nil {
		logx.Errorf("get user info id: %d err: %v", articleInfo.Article.AuthorId, err)
		return nil, err
	}
	publishTime := time.Unix(articleInfo.Article.PublishTime, 0).Format("2006-01-02 15:04:05")
	return &types.ArticleDetailResponse{
		ArticleId:    articleInfo.Article.Id,
		Title:        articleInfo.Article.Title,
		Content:      articleInfo.Article.Content,
		Description:  articleInfo.Article.Description,
		Cover:        articleInfo.Article.Cover,
		AuthorId:     articleInfo.Article.AuthorId,
		AuthorName:   userInfo.Username,
		AuthorAvatar: userInfo.Avatar,
		PublishTime:  publishTime,
		LikeNum:      articleInfo.Article.LikeCount,
	}, nil
}
