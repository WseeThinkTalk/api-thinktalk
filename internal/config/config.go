package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

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
