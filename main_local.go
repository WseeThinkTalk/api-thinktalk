//go:build local
// +build local

package main

import (
	"flag"
	"fmt"
	"net/http"

	"api-thinktalk/internal/config"
	"api-thinktalk/internal/handler"
	"api-thinktalk/internal/svc"
	"api-thinktalk/pkg/env"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func runLocal() {
	configFile := flag.String("f", "etc/api.yaml", "the config file")
	flag.Parse()

	env.LoadEnv()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(func(header http.Header) {
		header.Add("Access-Control-Allow-Headers", "x-token, Authorization, Content-Type")
	}, func(http.ResponseWriter) {}, "*"))
	defer server.Stop()

	// 健康检查
	server.AddRoutes([]rest.Route{
		{
			Method: http.MethodGet,
			Path:   "/health",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				httpx.OkJson(w, "ok")
			},
		},
	})

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting local thinktalk-api gateway at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
