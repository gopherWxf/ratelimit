package main

import (
	"fmt"
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

	fmt.Println("http://127.0.0.1:80")
	r.Run(":80")
}
