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
	//参数限流
	r.GET("/param", src.ParamLimiter(3, 1, "name")(test))
	//全局限流
	r.GET("/global", src.Limiter(3, 1)(test))
	//组合限流
	r.GET("/pg", src.ParamLimiter(3, 1, "name")(src.Limiter(1, 1)(test)))
	//IP限流
	r.GET("/ip", src.IPLimiter(3, 1)(test))

	fmt.Println("http://127.0.0.1:80")
	r.Run(":80")
}
