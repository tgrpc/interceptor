package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tb := NewTokenBucket(5)
	time.Sleep(time.Second * 1)
	tb.Limit()
	time.Sleep(time.Millisecond * 300)
	tb.Limit()
	time.Sleep(time.Second * 1)
	tb.Limit()
	time.Sleep(time.Millisecond * 300)
	tb.Limit()
	// <-make(chan bool, 1)
}
