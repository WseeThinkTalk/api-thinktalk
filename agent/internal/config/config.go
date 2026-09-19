package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth     struct {
		AccessSecret string
		AccessExpire int64
	}
	AgentRpc struct {
		Endpoints []string
		NonBlock  bool
		Timeout   int64
	}
}
