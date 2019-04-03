package webserver

import (
	"fmt"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/zekroTJA/slms/pkg/ratelimit"
	"github.com/zekroTJA/slms/pkg/timedmap"
)

const (
	cleanupInterval = 15 * time.Minute
	entryLifetime   = 1 * time.Hour
)

// A RateLimitManager maintains all
// rate limiters for each connection.
type RateLimitManager struct {
	limits *timedmap.TimedMap
}

// A RateLimitHandler provides a
// fasthttp-routing handler for
// per-route connection-based rate
// limiting.
type RateLimitHandler struct {
	routing.Handler

	limit time.Duration
	burst int
}

// NewRateLimitManager creates a new instance
// of RateLimitManager.
func NewRateLimitManager() *RateLimitManager {
	return &RateLimitManager{
		limits: timedmap.New(cleanupInterval),
	}
}

// NewRateLimitHandler creates a new per-route
// connection-based limiter handler.
func (rlm *RateLimitManager) NewRateLimitHandler(limit time.Duration, burst int) *RateLimitHandler {
	handler := func(ctx *routing.Context) error {
		ok, res := rlm.getLimiter(ctx.RemoteIP().String(), limit, burst).Reserve()

		ctx.Response.Header.Set("X-RateLimit-Limit", fmt.Sprintf("%d", res.Burst))
		ctx.Response.Header.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
		ctx.Response.Header.Set("X-RateLimit-Reset", fmt.Sprintf("%d", res.Reset.Unix()))

		if !ok {
			ctx.Abort()
			ctx.Response.Header.SetContentType("application/json")
			ctx.SetStatusCode(429)
			ctx.SetBodyString(
				"{\n  \"code\": 429,\n  \"message\": \"you are being rate limited\"\n}")
		}

		return nil
	}

	return &RateLimitHandler{
		limit:   limit,
		burst:   burst,
		Handler: handler,
	}
}

// getLimiter tries to get an existent limiter
// from the limiter map. If there is no limiter
// existent for this address, a new limiter
// will be created and added to the map.
func (rlm *RateLimitManager) getLimiter(addr string, limit time.Duration, burst int) *ratelimit.Limiter {
	var ok bool
	var limiter *ratelimit.Limiter

	if rlm.limits.Contains(addr) {
		limiter, ok = rlm.limits.GetValue(addr).(*ratelimit.Limiter)
		if !ok {
			limiter = rlm.createLimiter(addr, limit, burst)
		}
	} else {
		limiter = rlm.createLimiter(addr, limit, burst)
	}

	return limiter
}

// createLimiter creates a new limiter and
// adds it to the limiters map by the passed
// address.
func (rlm *RateLimitManager) createLimiter(addr string, limit time.Duration, burst int) *ratelimit.Limiter {
	limiter := ratelimit.NewLimiter(limit, burst)
	rlm.limits.Set(addr, limiter, entryLifetime)
	return limiter
}
