package ratelimit

import (
	"net/http"
	"time"
)

var DefaultHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(429)
}

func InitRateLimit(requests int, limitTime time.Duration, blockedHandler http.HandlerFunc) *RateLimiter {
	ratelimiter := &RateLimiter{
		requests:       requests,
		addresses:      make(map[string]Request),
		blockedHandler: blockedHandler,
		timeLimit:      limitTime,
		ValidateByURI:  true,
		IsUsingProxy:   false,
	}
	go ratelimiter.reduceTheLimits()
	return ratelimiter
}

func (rl *RateLimiter) RateLimit(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.RemoteAddr
		if rl.IsUsingProxy {
			remoteIP = r.Header.Get("REMOTE_ADDR")
		}
		if rl.ValidateByURI {
			remoteIP += r.RequestURI
		}
		if rl.exceededTheLimit(remoteIP) {			
			rl.blockedHandler.ServeHTTP(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func (rl *RateLimiter) reduceTheLimits() {
	rl.mux.Lock()
	defer rl.mux.Unlock()
	for key, value := range rl.addresses {
		if value.Time < time.Now().Unix()-int64(rl.timeLimit.Seconds()) {
			delete(rl.addresses, key)
		}
	}
	time.AfterFunc(rl.timeLimit, rl.reduceTheLimits)
}

func (rl *RateLimiter) exceededTheLimit(remoteIP string) bool {
	rl.mux.Lock()
	defer rl.mux.Unlock()
	req := rl.addresses[remoteIP]
	req.Count++
	req.Time = time.Now().Unix()
	rl.addresses[remoteIP] = req
	return req.Count > rl.requests
}
