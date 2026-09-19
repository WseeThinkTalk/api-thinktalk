// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

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
	Alipay struct {
		AppId           string
		PrivateKey      string
		AlipayPublicKey string
		NotifyURL       string
		ReturnURL       string
		IsProduction    bool
	}
	Redis       redis.RedisConf
	UserRpc     zrpc.RpcClientConf
	LikeRpc     zrpc.RpcClientConf
	FollowRpc   zrpc.RpcClientConf
	MessageRpc  zrpc.RpcClientConf
	ConcernedRpc zrpc.RpcClientConf
	MemberRpc   zrpc.RpcClientConf
	TagRpc      zrpc.RpcClientConf
	ReplyRpc    zrpc.RpcClientConf
	ArticleRpc  zrpc.RpcClientConf
	NotificationPusher struct {
		Brokers []string
		Topic   string
	}
}
