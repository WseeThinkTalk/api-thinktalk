package main

import (
	"flag"
	"fmt"

	agentapi "api-thinktalk/agent"
	appletapi "api-thinktalk/applet"
	articleapi "api-thinktalk/article"
	chatapi "api-thinktalk/chat"
	"api-thinktalk/pkg/env"
	"api-thinktalk/pkg/xcode"
	qaapi "api-thinktalk/qa"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
)

var configFile = flag.String("f", "etc/api.yaml", "the config file")

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
		RefreshAfter  int64
	}

	UserRpc    zrpc.RpcClientConf
	ContentRpc zrpc.RpcClientConf
	SocialRpc  zrpc.RpcClientConf
	AgentRpc   zrpc.RpcClientConf

	Redis    redis.RedisConf
	BizRedis redis.RedisConf

	MinIO struct {
		Endpoint        string
		AccessKeyID     string
		AccessKeySecret string
		BucketName      string
		Location        string
		UseSSL          bool
	}
	Alipay struct {
		AppId           string
		PrivateKey      string
		AlipayPublicKey string
		NotifyURL       string
		ReturnURL       string
		IsProduction    bool
	}
	NotificationPusher struct {
		Brokers []string
		Topic   string
	}
}

func main() {
	flag.Parse()

	env.LoadEnv()

	var c Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 1. Applet Routes (User, Member, Follow, Like, Concerned, Message, Tag, Reply)
	appletConf := appletapi.Config{
		RestConf:           c.RestConf,
		Auth:               c.Auth,
		UserRpc:            c.UserRpc,
		LikeRpc:            c.SocialRpc,
		FollowRpc:          c.UserRpc,
		MessageRpc:         c.SocialRpc,
		ConcernedRpc:       c.SocialRpc,
		MemberRpc:          c.UserRpc,
		TagRpc:             c.ContentRpc,
		ReplyRpc:           c.SocialRpc,
		ArticleRpc:         c.ContentRpc,
		Redis:              c.Redis,
		Alipay:             c.Alipay,
		NotificationPusher: c.NotificationPusher,
	}
	appletapi.RegisterRoutes(server, appletConf)

	// 2. Article Routes
	articleConf := articleapi.Config{
		RestConf: c.RestConf,
		Auth: struct {
			AccessSecret string
			AccessExpire int64
		}{
			AccessSecret: c.Auth.AccessSecret,
			AccessExpire: c.Auth.AccessExpire,
		},
		MinIO:      c.MinIO,
		ArticleRPC: c.ContentRpc,
		UserRPC:    c.UserRpc,
	}
	articleapi.RegisterRoutes(server, articleConf)

	// 3. Chat Routes (including WebSocket)
	chatConf := chatapi.Config{
		RestConf: c.RestConf,
		Auth: struct {
			AccessSecret string
			AccessExpire int64
		}{
			AccessSecret: c.Auth.AccessSecret,
			AccessExpire: c.Auth.AccessExpire,
		},
		ChatRpc:  c.SocialRpc,
		BizRedis: c.BizRedis,
		UserRpc:  c.UserRpc,
	}
	chatapi.RegisterRoutes(server, chatConf)

	// 4. QA Routes
	qaConf := qaapi.Config{
		RestConf: c.RestConf,
		Auth: struct {
			AccessSecret string
			AccessExpire int64
		}{
			AccessSecret: c.Auth.AccessSecret,
			AccessExpire: c.Auth.AccessExpire,
		},
		QaRpc: c.ContentRpc,
	}
	qaapi.RegisterRoutes(server, qaConf)

	// 5. Agent Routes
	agentConf := agentapi.Config{
		RestConf: c.RestConf,
		Auth: struct {
			AccessSecret string
			AccessExpire int64
		}{
			AccessSecret: c.Auth.AccessSecret,
			AccessExpire: c.Auth.AccessExpire,
		},
		AgentRpc: struct {
			Endpoints []string
			NonBlock  bool
			Timeout   int64
		}{
			Endpoints: c.AgentRpc.Endpoints,
			NonBlock:  c.AgentRpc.NonBlock,
			Timeout:   c.AgentRpc.Timeout,
		},
	}
	agentapi.RegisterRoutes(server, agentConf)

	httpx.SetErrorHandler(xcode.ErrHandler)

	fmt.Printf("Starting unified thinktalk-api gateway at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
