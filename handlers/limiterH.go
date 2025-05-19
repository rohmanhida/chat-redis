package handlers

import (
	"sync"

	"golang.org/x/time/rate"
)

type client struct {
	limiter *rate.Limiter
}

var (
	visitors = make(map[string]*client)
	mu       sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(15, 5) // 1 msg/sec, burst of 5
		visitors[ip] = &client{limiter}
		return limiter
	}
	return v.limiter
}
