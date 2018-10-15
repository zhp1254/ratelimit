package api

import (
	"github.com/zhp1254/ratelimit/leakybucket"
	"github.com/zhp1254/ratelimit/simpleratelimit"
	"time"
)

type Limiter interface {
  	Allow() bool
}

func NewLimiter(name string, capacity uint, per time.Duration) Limiter{
	if name == "leaky"{
		return leakybucket.NewLimiter(capacity, per)
	}else if name == "token"{
		return simpleratelimit.NewLimiter(capacity, per)
	}
	panic("unknow limiter")
}