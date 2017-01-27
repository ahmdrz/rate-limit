package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

type Request struct {
	Time  int64
	Count int
}

type RateLimiter struct {
	IsUsingProxy   bool
	ValidateByURI  bool
	requests       int
	addresses      map[string]Request
	blockedHandler http.HandlerFunc
	mux            sync.Mutex
	timeLimit      time.Duration
}

type HandlerLimiter struct {
	RateLimiter
	http.Handler
}
