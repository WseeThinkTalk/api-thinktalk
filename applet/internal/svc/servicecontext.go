package svc

import (
	"strings"

	"api-thinktalk/applet/internal/config"
	"api-thinktalk/applet/internal/middleware"
	"api-thinktalk/client/article/article"
	"api-thinktalk/client/concerned/concerned"
	"api-thinktalk/client/follow/follow"
	like "api-thinktalk/client/like/like"
	"api-thinktalk/client/member/member"
	"api-thinktalk/client/message/message"
	"api-thinktalk/client/reply/reply"
	"api-thinktalk/client/tag/tag"
	"api-thinktalk/client/user/user"
	"api-thinktalk/pkg/interceptors"

	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	RDB          *redis.Redis
	UserRPC      user.User
	LikeRPC      like.Like
	FollowRPC    follow.Follow
	MessageRPC   message.Message
	ConcernedRPC concerned.Concerned
	MemberRPC    member.Member
	TagRPC       tag.Tag
	ReplyRPC     reply.Reply
	ArticleRPC   article.Article
	AdminAuth    rest.Middleware
	AlipayClient *alipay.Client
	NotificationPusher *kq.Pusher
}

func formatPEMKey(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return s
	}
	if strings.Contains(s, "-----BEGIN") {
		return s
	}
	return "-----BEGIN PRIVATE KEY-----\n" + s + "\n-----END PRIVATE KEY-----"
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb, _ := redis.NewRedis(c.Redis)
	userRPC := zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	likeRPC := zrpc.MustNewClient(c.LikeRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	followRPC := zrpc.MustNewClient(c.FollowRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	messageRPC := zrpc.MustNewClient(c.MessageRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	concernedRPC := zrpc.MustNewClient(c.ConcernedRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	memberRPC := zrpc.MustNewClient(c.MemberRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	tagRPC := zrpc.MustNewClient(c.TagRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	replyRPC := zrpc.MustNewClient(c.ReplyRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	articleRPC := zrpc.MustNewClient(c.ArticleRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	userRPCModel := user.NewUser(userRPC)

	logx.Infof("[Alipay Config] AppId='%s', PrivateKeyLen=%d, PublicKeyLen=%d, IsProduction=%t, NotifyURL='%s'",
		c.Alipay.AppId, len(c.Alipay.PrivateKey), len(c.Alipay.AlipayPublicKey), c.Alipay.IsProduction, c.Alipay.NotifyURL)

	var alipayClient *alipay.Client
	if c.Alipay.AppId != "" && c.Alipay.PrivateKey != "" {
		privateKey := formatPEMKey(c.Alipay.PrivateKey)
		var err error
		alipayClient, err = alipay.New(c.Alipay.AppId, privateKey, c.Alipay.IsProduction)
		if err != nil {
			logx.Errorf("[Alipay] Init client failed: %v. Payment disabled, check ALIPAY_PRIVATE_KEY format.", err)
		} else if c.Alipay.AlipayPublicKey != "" {
			err = alipayClient.LoadAliPayPublicKey(c.Alipay.AlipayPublicKey)
			if err != nil {
				logx.Errorf("[Alipay] Load public key failed: %v", err)
				alipayClient = nil
			}
		}
	}

	if alipayClient != nil {
		logx.Info("[Alipay] Client initialized successfully")
	} else {
		logx.Info("[Alipay] Client NOT initialized — running without payment (Mock mode)")
	}

	return &ServiceContext{
		Config:             c,
		RDB:                rdb,
		UserRPC:            userRPCModel,
		LikeRPC:            like.NewLike(likeRPC),
		FollowRPC:          follow.NewFollow(followRPC),
		MessageRPC:         message.NewMessage(messageRPC),
		ConcernedRPC:       concerned.NewConcerned(concernedRPC),
		MemberRPC:          member.NewMember(memberRPC),
		TagRPC:             tag.NewTag(tagRPC),
		ReplyRPC:           reply.NewReply(replyRPC),
		ArticleRPC:         article.NewArticle(articleRPC),
		AdminAuth:          middleware.NewAdminAuthMiddleware(userRPCModel).Handle,
		AlipayClient:       alipayClient,
		NotificationPusher: kq.NewPusher(c.NotificationPusher.Brokers, c.NotificationPusher.Topic),
	}
}
