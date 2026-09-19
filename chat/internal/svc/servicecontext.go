package svc

import (
	"api-thinktalk/chat/internal/config"
	"api-thinktalk/chat/internal/hub"
	"api-thinktalk/client/chat/chat"
	"api-thinktalk/client/user/user"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	Chat     chat.Chat
	UserRPC  user.User
	Redis    *redis.Redis
	Hub      *hub.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(redis.RedisConf{
		Host:        c.BizRedis.Host,
		Pass:        c.BizRedis.Pass,
		Type:        c.BizRedis.Type,
		PingTimeout: 10000000000,
	})

	return &ServiceContext{
		Config:  c,
		Chat:    chat.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		UserRPC: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Redis:   rds,
		Hub:     hub.NewHub(),
	}
}
