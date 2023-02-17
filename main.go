package main

import (
	"github.com/gin-gonic/gin"
	"ratelimit/src"
)

func main() {
	ratelimter := src.NewBucket(10)

	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		if ratelimter.IsAccept() {
			ctx.JSON(200, gin.H{"message": "ok"})
		} else {
			ctx.AbortWithStatusJSON(400, gin.H{"message": "rate limit"})
		}
	})
	r.Run(":80")
}
