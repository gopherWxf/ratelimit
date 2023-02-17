package src

import "github.com/gin-gonic/gin"

// 限流装饰器
func Limiter(cap int) func(handler gin.HandlerFunc) gin.HandlerFunc {
	limiter := NewBucket(cap)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			if limiter.IsAccept() {
				handler(ctx)
			} else {
				ctx.AbortWithStatusJSON(429, gin.H{"message": "too many requests"})
			}
		}
	}
}