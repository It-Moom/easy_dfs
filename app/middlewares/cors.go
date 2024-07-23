/*
 * @PackageName: middlewares
 * @FileName: cors.go
 * @Description: 跨域中间件
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:45:17
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:45:17
 */

package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// 跨域处理函数
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin") // 请求头部
		if origin != "" {
			// 可将将 * 替换为指定的域名
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		ctx.Next()
	}
}
