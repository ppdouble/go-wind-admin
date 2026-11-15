package rpc

import (
    "github.com/go-kratos/aegis/ratelimit"
    kratoslimit "github.com/go-kratos/kratos/v3/middleware/ratelimit"
)

// aegisLimiter preserves the configured BBR limiter across v3 callback types.
type aegisLimiter struct{ ratelimit.Limiter }

func (l aegisLimiter) Allow() (kratoslimit.DoneFunc, error) {
    done, err := l.Limiter.Allow()
    if err != nil { return nil, err }
    return func(info kratoslimit.DoneInfo) { done(ratelimit.DoneInfo{Err: info.Err}) }, nil
}
