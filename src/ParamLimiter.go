package src

import (
	"github.com/gin-gonic/gin"
)

// ParamLimiter 基于参数限流 http://127.0.0.1/?name=wxf   name就是key参数
func ParamLimiter(cap int64, rate int64, key string) func(handler gin.HandlerFunc) gin.HandlerFunc {
	limiter := NewBucket(cap, rate)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			//有参数，做限流
			if ctx.Query(key) != "" {
				if limiter.IsAccept() {
					handler(ctx)
				} else {
					ctx.AbortWithStatusJSON(429, gin.H{"message": "too many requests-param"})
				}
			} else {
				handler(ctx)
			}
		}
	}
}
