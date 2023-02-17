package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ratelimit/src"
)

func test(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "ok"})
}
func main() {
	r := gin.New()
	r.GET("/", src.Limiter(10)(test))

	fmt.Println("http://127.0.0.1:80")
	r.Run(":80")
}
