package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/internal/config"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func init() {
	go cleanupClients()
}

// RateLimiter returns a gin.HandlerFunc that limits requests per IP
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := config.HTTP.RateLimitPerMinute
		if limit <= 0 {
			c.Next()
			return
		}

		ip := c.ClientIP()
		mu.Lock()
		v, exists := clients[ip]
		if !exists {
			// rate.Limit is requests per second.
			// We want limit per minute, so we divide by 60.
			limiter := rate.NewLimiter(rate.Limit(float64(limit)/60.0), limit)
			clients[ip] = &client{limiter: limiter, lastSeen: time.Now()}
			v = clients[ip]
		}
		v.lastSeen = time.Now()
		mu.Unlock()

		if !v.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// cleanupClients removes old client entries from the map to prevent memory leaks
func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
