// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"api-thinktalk/article/internal/config"
	"api-thinktalk/article/internal/middleware"
	"api-thinktalk/client/article/article"
	"api-thinktalk/client/user/user"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRPC    user.User
	ArticleRPC article.Article
	MinIO      *minio.Client
	AdminAuth  rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	minioClient, err := minio.New(c.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MinIO.AccessKeyID, c.MinIO.AccessKeySecret, ""),
		Secure: c.MinIO.UseSSL,
	})
	if err != nil {
		panic("minio 连接失败: " + err.Error())
	}

	userRPC := user.NewUser(zrpc.MustNewClient(c.UserRPC))
	return &ServiceContext{
		Config:     c,
		UserRPC:    userRPC,
		ArticleRPC: article.NewArticle(zrpc.MustNewClient(c.ArticleRPC)),
		MinIO:      minioClient,
		AdminAuth:  middleware.NewAdminAuthMiddleware(userRPC).Handle,
	}
}
