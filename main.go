package main

import (
	"flag"
	"fmt"
	"net/http"

	"api-thinktalk/internal/config"
	"api-thinktalk/internal/handler"
	"api-thinktalk/internal/svc"
	"api-thinktalk/pkg/lib/etcdx"
	"api-thinktalk/pkg/lib/zapx"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

func runRemoteConfig() *config.Config {
	var c config.Config
	etcdx.MustLoadRemoteConfig("/thinktalk/config/api", &c)
	return &c
}

func main() {
	flag.Parse()

	// 从 Etcd 配置中心拉取远程配置 (Fail-Fast)
	c := runRemoteConfig()
	if c == nil {
		return
	}

	// init logger
	writer, err := zapx.NewZapWriter()
	if err == nil {
		logx.SetWriter(writer)
	}

	// 设置请求体最大限制 100MB
	c.RestConf.MaxBytes = 100 << 20

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(func(header http.Header) {
		header.Add("Access-Control-Allow-Headers", "x-token, Authorization, Content-Type")
	}, func(http.ResponseWriter) {}, "*"))
	defer server.Stop()

	ctx := svc.NewServiceContext(*c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting unified thinktalk-api gateway at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
