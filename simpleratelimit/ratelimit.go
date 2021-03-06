package simpleratelimit

import (
	"sync/atomic"
	"time"
)

// RateLimiter 限速器
type RateLimiter struct {
	rate      uint64  //单位时间可用数量
	allowance uint64  //当前可用数量 剩余的
	max       uint64  //最大数量
	unit      uint64  //限制最小时间单位 默认 1秒
	lastCheck uint64
}

// New 创建RateLimiter实例
func NewLimiter(capacity uint, per time.Duration) *RateLimiter {
	nano := uint64(per)
	if nano < 1 {
		nano = uint64(time.Second)
	}
	if capacity < 1 {
		capacity = 1
	}

	return &RateLimiter{
		rate:      uint64(capacity),
		allowance: uint64(capacity) * nano,
		max:       uint64(capacity) * nano,
		unit:      nano,

		lastCheck: unixNano(),
	}
}

//alias Limit()
func (rl *RateLimiter) Allow() bool {
	return !rl.Limit()
}

// Limit 判断是否超过限制
func (rl *RateLimiter) Limit() bool {
	now := unixNano()
	// 计算上一次调用到现在过了多少纳秒
	passed := now - atomic.SwapUint64(&rl.lastCheck, now)

	rate := atomic.LoadUint64(&rl.rate)
	current := atomic.AddUint64(&rl.allowance, passed*rate)

	if max := atomic.LoadUint64(&rl.max); current > max {
		atomic.AddUint64(&rl.allowance, max-current)
		current = max
	}

	if current < rl.unit {
		return true
	}

	// 没有超过限额
	atomic.AddUint64(&rl.allowance, -rl.unit)
	return false
}

// UpdateRate 更新速率值
func (rl *RateLimiter) UpdateRate(rate int) {
	atomic.StoreUint64(&rl.rate, uint64(rate))
	atomic.StoreUint64(&rl.max, uint64(rate)*rl.unit)
}

// Undo 重置上一次调用Limit()，返回没有使用过的限额
func (rl *RateLimiter) Undo() {
	current := atomic.AddUint64(&rl.allowance, rl.unit)

	if max := atomic.LoadUint64(&rl.max); current > max {
		atomic.AddUint64(&rl.allowance, max-current)
	}
}

// unixNano 当前时间（纳秒）
func unixNano() uint64 {
	return uint64(time.Now().UnixNano())
}
