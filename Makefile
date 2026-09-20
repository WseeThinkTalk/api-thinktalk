run:
	go run main.go

file:
	goctl api go --api ./doc/api.api -dir . --style=go_zero

format:
	goctl api format -dir ./doc/user/user.api
	goctl api format -dir ./doc/article/article.api
	goctl api format -dir ./doc/chat/chat.api
	goctl api format -dir ./doc/agent/agent.api
	goctl api format -dir ./doc/qa/qa.api
	goctl api format -dir ./doc/social/social.api
	goctl api format -dir ./doc/member/member.api
	goctl api format -dir ./doc/tag/tag.api

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api-thinktalk main.go
