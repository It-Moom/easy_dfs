/*
 * @PackageName: middlewares
 * @FileName: file_url_check.go
 * @Description: 文件URL检查中间件
 * @Author: gabbymrh
 * @Date: 2024-07-20 10:15:52
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-20 10:15:52
 */

package middlewares

import (
	"easy_dfs/app/enum/response_code"
	"easy_dfs/app/enum/system_default"
	"easy_dfs/app/services"
	"easy_dfs/pkg/http/http_response"
	"errors"
	"github.com/gin-gonic/gin"
	"net/url"
	"strings"
)

// 文件URL检查
func FileUrlCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取URL
		fileUrl := ctx.Request.URL
		// 如果访问URL包含storage则校验bucket
		if strings.Contains(fileUrl.String(), system_default.STORAGE_PATH) {
			// fmt.Println("fileUrl:", fileUrl)
			// 获取bucket
			queryParams, bcerr := url.ParseQuery(ctx.Request.URL.RawQuery)
			if bcerr != nil {
				http_response.Response(ctx, response_code.PARAM_ERROR, false, "操作失败", nil, bcerr)
				return
			}
			// 获取bucket
			bucket := queryParams.Get("bucket")
			if bucket == "" {
				http_response.Response(ctx, response_code.PARAM_ERROR, false, "操作失败", nil, errors.New("bucket不能为空"))
				return
			}
			bs := new(services.BucketService)
			// 校验bucket
			_, fberr := bs.FindBucketInfo(bucket)
			if fberr != nil {
				http_response.Response(ctx, response_code.REQUEST_FAILS, false, "操作失败", nil, fberr)
				return
			}
		}

		ctx.Next()

	}
}
