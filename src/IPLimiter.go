package src

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LimiterCache struct {
	data sync.Map //key->ip+port val->bucket
}

var IpCache2 *GCache
var IpCache *LimiterCache

func init() {
	IpCache = &LimiterCache{}
	IpCache2 = NewGCache(WithMaxSize(10000))
}

// IPLimiter ip限流
func IPLimiter(cap int64, rate int64) func(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			ip := ctx.ClientIP()
			var limiter *Bucket
			//if v, ok := IpCache.data.Load(ip); ok {
			//	limiter = v.(*Bucket)
			//} else {
			//	limiter = NewBucket(cap, rate)
			//	IpCache.data.Store(ip, limiter)
			//}
			if v := IpCache2.Get(ip); v != nil {
				limiter = v.(*Bucket)
			} else {
				fmt.Println("form cache")
				limiter = NewBucket(cap, rate)
				IpCache2.Set(ip, limiter, time.Second*5)
			}

			if limiter.IsAccept() {
				handler(ctx)
			} else {
				ctx.AbortWithStatusJSON(429, gin.H{"message": "too many requests-ip"})
			}
		}
	}
}

// 代理
func ClientIP(r *http.Request) string {
	ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
