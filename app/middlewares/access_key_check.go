/*
 * @PackageName: middlewares
 * @FileName: access_key_check.go
 * @Description: 访问密钥检查中间件
 * @Author: gabbymrh
 * @Date: 2024-07-19 18:11:15
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-19 18:11:15
 */

package middlewares

import (
	"easy_dfs/app/enum/response_code"
	"easy_dfs/app/services"
	"easy_dfs/pkg/http/http_response"
	"errors"
	"github.com/gin-gonic/gin"
)

// 校验访问密钥
func AccessKeyCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取访问密钥
		accessKey := ctx.GetHeader("X-Access-Key")
		secretKey := ctx.GetHeader("X-Secret-Key")
		if accessKey == "" || secretKey == "" {
			http_response.Response(ctx, response_code.TOKEN_INVALID, false, "操作失败", nil, errors.New("密钥不能为空"))
			return
		}

		// 校验访问密钥
		accessKeyService := new(services.AccessKeyService)
		checkKey := accessKeyService.CheckAccessKey(accessKey, secretKey)
		if !checkKey {
			http_response.Response(ctx, response_code.TOKEN_INVALID, false, "操作失败", nil, errors.New("密钥无效"))
			return
		} else {
			ctx.Next()
		}

	}
}
