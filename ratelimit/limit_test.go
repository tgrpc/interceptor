package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(20)
	for i := 0; i < 30; i++ {
		time.Sleep(time.Millisecond * 100)
		if !tb.Limit() {
			fmt.Println("limit fail", tb.Length())
			continue
		}
		fmt.Println("limit ", tb.Length())
	}
	// <-make(chan bool, 1)
}

// <30 ns/op; 不用原子操作，60 ns/op
func BenchmarkPut(b *testing.B) {
	tb := NewTokenBucket(1000)
	// 不开启loopPut，测试put的性能
	for i := 0; i < b.N; i++ {
		tb.put()
	}
}

// 25 ns/op; 不用原子操作，60+ ns/op
func BenchmarkLimit(b *testing.B) {
	tb := NewTokenBucket(1000)
	for i := 0; i < b.N; i++ {
		tb.Limit()
	}
}
