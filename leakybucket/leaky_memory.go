package leakybucket

import (
	"sync"
	"time"
)

type LeakyLimiter struct {
	capacity  uint     //容量
	remaining uint     //剩余
	reset     time.Time
	rate      time.Duration //限制时间
	mutex     sync.Mutex
}

// Capacity return
func (b *LeakyLimiter) Capacity() uint {
	return b.capacity
}

// Remaining space in the LeakyLimiter.
func (b *LeakyLimiter) Remaining() uint {
	return b.remaining
}

// Reset returns when the LeakyLimiter will be drained.
func (b *LeakyLimiter) Reset() time.Time {
	b.remaining = b.capacity
	return b.reset
}

// Add to the LeakyLimiter.
func (b *LeakyLimiter) Add(amount uint) (BucketState, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if time.Now().After(b.reset) {
		b.reset = time.Now().Add(b.rate)
		b.remaining = b.capacity
	}
	if amount > b.remaining {
		return BucketState{Capacity: b.capacity, Remaining: b.remaining, Reset: b.reset}, ErrorFull
	}
	b.remaining -= amount
	return BucketState{Capacity: b.capacity, Remaining: b.remaining, Reset: b.reset}, nil
}

//alias Add(1) 占用一个token
func (b *LeakyLimiter) Allow() bool{
	_, err := b.Add(1)
	return err == nil
}


/**
 * remaining 初始剩余量
 * capacity  问题
 * rate      限制时间
 */
func NewLimiter(capacity uint, per time.Duration) *LeakyLimiter {
	return &LeakyLimiter{
		capacity:  capacity,
		remaining: capacity,
		reset:     time.Now().Add(per),
		rate:      per,
	}
}

