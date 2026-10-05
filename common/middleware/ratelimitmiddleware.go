package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Period int // 周期（秒）
	Quota  int // 配额（次）
}

// RateLimitMiddleware 路由级自适应限流中间件
type RateLimitMiddleware struct {
	cfg     RateLimitConfig
	redis   *redis.Redis
	limiter *limit.PeriodLimit
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware(cfg RateLimitConfig, rds *redis.Redis) *RateLimitMiddleware {
	if cfg.Period <= 0 {
		cfg.Period = 1
	}
	if cfg.Quota <= 0 {
		cfg.Quota = 100 // 默认每秒 100 次
	}

	var limiter *limit.PeriodLimit
	if rds != nil {
		limiter = limit.NewPeriodLimit(cfg.Period, cfg.Quota, rds, "rate_limit:")
	}

	return &RateLimitMiddleware{
		cfg:     cfg,
		redis:   rds,
		limiter: limiter,
	}
}

// Handle 拦截处理
func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.limiter == nil {
			next(w, r)
			return
		}

		// 按 IP 与路径联合限流
		clientIP := getClientIP(r)
		key := fmt.Sprintf("%s:%s", r.URL.Path, clientIP)

		code, err := m.limiter.Take(key)
		if err != nil {
			// 限流器异常平滑放行，不阻断正常业务
			next(w, r)
			return
		}

		if code == limit.OverQuota {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"code":429,"msg":"请求过于频繁，请稍后重试"}`))
			return
		}

		next(w, r)
	}
}

func getClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
