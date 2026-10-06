// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"
	"time"

	"api-thinktalk/internal/svc"
	"api-thinktalk/internal/types"
	article "api-thinktalk/client/article/pb"
	user "api-thinktalk/client/user/service"

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
	resp = new(types.ArticleDetailResponse)

	articleInfo, err := l.svcCtx.ArticleRPC.ArticleDetail(l.ctx, &article.ArticleDetailRequest{
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("get article detail id: %d err: %v", req.ArticleId, err)
		return nil, err
	}
	if articleInfo == nil || articleInfo.Data == nil {
		return nil, nil
	}
	artData := articleInfo.Data
	userInfo, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: artData.AuthorId,
	})
	if err != nil {
		logx.Errorf("get user info id: %d err: %v", artData.AuthorId, err)
		return nil, err
	}

	authorName := ""
	authorAvatar := ""
	if userInfo != nil && userInfo.Data != nil {
		authorName = userInfo.Data.Username
		authorAvatar = userInfo.Data.Avatar
	}

	publishTime := time.Unix(artData.PublishTime, 0).Format("2006-01-02 15:04:05")
	resp.ArticleId = artData.Id
	resp.Title = artData.Title
	resp.Content = artData.Content
	resp.Description = artData.Description
	resp.Cover = artData.Cover
	resp.AuthorId = artData.AuthorId
	resp.AuthorName = authorName
	resp.AuthorAvatar = authorAvatar
	resp.PublishTime = publishTime
	resp.LikeNum = artData.LikeCount
	return resp, nil
}
