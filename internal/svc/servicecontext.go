package svc

import (
	"strings"

	agent "api-thinktalk/client/agent/pb"
	article "api-thinktalk/client/article/pb"
	chat "api-thinktalk/client/chat/pb"
	concerned "api-thinktalk/client/concerned/pb"
	follow "api-thinktalk/client/follow/pb"
	like "api-thinktalk/client/like/service"
	member "api-thinktalk/client/member/pb"
	message "api-thinktalk/client/message/pb"
	qa "api-thinktalk/client/qa/pb"
	reply "api-thinktalk/client/reply/pb"
	tag "api-thinktalk/client/tag/pb"
	user "api-thinktalk/client/user/service"
	"api-thinktalk/internal/config"
	"api-thinktalk/common/hub"
	"api-thinktalk/common/middleware"
	"api-thinktalk/pkg/interceptors"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	RDB                *redis.Redis
	Redis              *redis.Redis // Alias for chat logic
	BizRedis           *redis.Redis
	UserRPC            user.UserClient
	LikeRPC            like.LikeClient
	FollowRPC          follow.FollowClient
	MessageRPC         message.MessageClient
	ConcernedRPC       concerned.ConcernedClient
	MemberRPC          member.MemberClient
	TagRPC             tag.TagClient
	ReplyRPC           reply.ReplyClient
	ArticleRPC         article.ArticleClient
	Chat               chat.ChatClient
	ChatRPC            chat.ChatClient
	QaRPC              qa.QAClient
	AgentClient        agent.AgentClient
	MinIO              *minio.Client
	AdminAuth          rest.Middleware
	RateLimit          rest.Middleware
	AlipayClient       *alipay.Client
	NotificationPusher *kq.Pusher
	Hub                *hub.Hub
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
	bizRedis := redis.MustNewRedis(redis.RedisConf{
		Host:        c.BizRedis.Host,
		Pass:        c.BizRedis.Pass,
		Type:        c.BizRedis.Type,
		PingTimeout: 10000000000,
	})

	// 4 Consolidated RPC Clients
	userClient := zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	contentClient := zrpc.MustNewClient(c.ContentRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	socialClient := zrpc.MustNewClient(c.SocialRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	agentClient := zrpc.MustNewClient(c.AgentRpc)

	userRPCModel := user.NewUserClient(userClient.Conn())
	articleRPCModel := article.NewArticleClient(contentClient.Conn())
	tagRPCModel := tag.NewTagClient(contentClient.Conn())
	qaRPCModel := qa.NewQAClient(contentClient.Conn())

	likeRPCModel := like.NewLikeClient(socialClient.Conn())
	followRPCModel := follow.NewFollowClient(userClient.Conn())
	messageRPCModel := message.NewMessageClient(socialClient.Conn())
	concernedRPCModel := concerned.NewConcernedClient(socialClient.Conn())
	memberRPCModel := member.NewMemberClient(userClient.Conn())
	replyRPCModel := reply.NewReplyClient(socialClient.Conn())
	chatRPCModel := chat.NewChatClient(socialClient.Conn())
	agentClientModel := agent.NewAgentClient(agentClient.Conn())

	// MinIO
	var minioClient *minio.Client
	if c.MinIO.Endpoint != "" {
		var err error
		minioClient, err = minio.New(c.MinIO.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(c.MinIO.AccessKeyID, c.MinIO.AccessKeySecret, ""),
			Secure: c.MinIO.UseSSL,
		})
		if err != nil {
			logx.Errorf("[MinIO] Connect failed: %v", err)
		}
	}

	// Alipay
	var alipayClient *alipay.Client
	if c.Alipay.AppId != "" && c.Alipay.PrivateKey != "" {
		privateKey := formatPEMKey(c.Alipay.PrivateKey)
		var err error
		alipayClient, err = alipay.New(c.Alipay.AppId, privateKey, c.Alipay.IsProduction)
		if err != nil {
			logx.Errorf("[Alipay] Init client failed: %v", err)
		} else if c.Alipay.AlipayPublicKey != "" {
			err = alipayClient.LoadAliPayPublicKey(c.Alipay.AlipayPublicKey)
			if err != nil {
				logx.Errorf("[Alipay] Load public key failed: %v", err)
				alipayClient = nil
			}
		}
	}

	var pusher *kq.Pusher
	if len(c.NotificationPusher.Brokers) > 0 && c.NotificationPusher.Topic != "" {
		pusher = kq.NewPusher(c.NotificationPusher.Brokers, c.NotificationPusher.Topic)
	}

	hubInstance := hub.NewHub()

	return &ServiceContext{
		Config:             c,
		RDB:                rdb,
		Redis:              bizRedis,
		BizRedis:           bizRedis,
		UserRPC:            userRPCModel,
		LikeRPC:            likeRPCModel,
		FollowRPC:          followRPCModel,
		MessageRPC:         messageRPCModel,
		ConcernedRPC:       concernedRPCModel,
		MemberRPC:          memberRPCModel,
		TagRPC:             tagRPCModel,
		ReplyRPC:           replyRPCModel,
		ArticleRPC:         articleRPCModel,
		Chat:               chatRPCModel,
		ChatRPC:            chatRPCModel,
		QaRPC:              qaRPCModel,
		AgentClient:        agentClientModel,
		MinIO:              minioClient,
		AdminAuth:          middleware.NewAdminAuthMiddleware(userRPCModel).Handle,
		RateLimit:          middleware.NewRateLimitMiddleware(middleware.RateLimitConfig{Period: 1, Quota: 100}, bizRedis).Handle,
		AlipayClient:       alipayClient,
		NotificationPusher: pusher,
		Hub:                hubInstance,
	}
}
