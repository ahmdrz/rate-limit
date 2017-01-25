package ratelimit

import (
	"net/http"
	"time"
)

var IsUsingProxy = false

var DefaultHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(429)
}

func InitRateLimit(requests int, limitTime time.Duration, blockedHandler http.HandlerFunc) *RateLimiter {
	ratelimiter := &RateLimiter{
		requests:       requests,
		addresses:      make(map[string]Request),
		blockedHandler: blockedHandler,
		timeLimit:      limitTime,
	}
	go ratelimiter.reduceTheLimits()
	return ratelimiter
}

func (rl *RateLimiter) RateLimit(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if IsUsingProxy {
			if rl.exceededTheLimit(r.RemoteAddr + r.RequestURI) {
				rl.blockedHandler.ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
		} else {
			if rl.exceededTheLimit(r.Header.Get("REMOTE_ADDR") + r.RequestURI) {
				rl.blockedHandler.ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
		}
	}
}

func (r *RateLimiter) reduceTheLimits() {
	r.mux.Lock()
	defer r.mux.Unlock()
	for key, value := range r.addresses {
		if value.Time < time.Now().Unix()-int64(r.timeLimit.Seconds()) {
			delete(r.addresses, key)
		}
	}
	time.AfterFunc(r.timeLimit, r.reduceTheLimits)
}

func (r *RateLimiter) exceededTheLimit(remoteIP string) bool {
	r.mux.Lock()
	defer r.mux.Unlock()
	req := r.addresses[remoteIP]
	req.Count++
	req.Time = time.Now().Unix()
	r.addresses[remoteIP] = req
	return req.Count > r.requests
}
