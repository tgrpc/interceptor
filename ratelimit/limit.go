package ratelimit

import (
	"sync"
	"sync/atomic"
	"time"
)

// 环形队列实现令牌桶限流
type TokenBucket struct {
	Size       int32
	head, rear int32
	sync.Mutex
}

func NewTokenBucket(size int32) *TokenBucket {
	if size <= 0 {
		panic("size should be bigger than 0")
	}
	tb := &TokenBucket{
		Size: size + 1,
	}
	go tb.loopPut()
	return tb
}

func (tb *TokenBucket) next(cur int32) int32 {
	return (cur + 1) % tb.Size
}

func (tb *TokenBucket) Length() int32 {
	return (tb.rear + tb.Size - tb.head) % tb.Size
}

func (tb *TokenBucket) Limit() bool {
	head := tb.head
	// 无令牌可用
	if atomic.CompareAndSwapInt32(&head, tb.rear, head) {
		return false
	}

	tb.Lock()
	defer tb.Unlock()
	// 无令牌可用
	if tb.head == tb.rear {
		return false
	}
	tb.head = tb.next(tb.head)
	return true
}

// 将令牌放入桶中
func (tb *TokenBucket) put() bool {
	next1 := tb.next(tb.rear)
	// 桶已满
	if atomic.CompareAndSwapInt32(&next1, tb.head, next1) {
		return false
	}

	// lock
	tb.Lock()
	defer tb.Unlock()
	// 桶已满，不需要放入
	next := tb.next(tb.rear)
	if next == tb.head {
		return false
	}
	tb.rear = next
	return true
}

func (tb *TokenBucket) loopPut() {
	dur := time.Duration(int64(time.Second) / int64(tb.Size-1))
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-ticker.C:
			tb.put()
		}
	}
}
