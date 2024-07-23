/*
 * @PackageName: middlewares
 * @FileName: recovery.go
 * @Description: 异常记录中间件
 * @Author: gabbymrh
 * @Date: 2024-07-23 14:03:50
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-23 14:03:50
 */

package middlewares

import (
	"easy_dfs/app/enum/response_code"
	"easy_dfs/pkg/http/http_response"
	"easy_dfs/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// Recovery 使用 zap.Error() 来记录 Panic 和 call stack
func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				// 获取用户的请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, true)

				// 链接中断，客户端中断连接为正常行为，不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// 链接中断的情况
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					// 链接已断开，无法写状态码
					return
				}

				// 如果不是链接中断，就开始记录堆栈信息
				logger.Error("recovery from panic",
					zap.Time("time", time.Now()),               // 记录时间
					zap.Any("error", err),                      // 记录错误信息
					zap.String("request", string(httpRequest)), // 请求信息
					zap.Stack("stacktrace"),                    // 调用堆栈信息
				)

				// 返回服务器错误信息
				http_response.Response(c, response_code.REQUEST_FAILS, false, "服务器内部错误", nil, nil)
				return
			}
		}()
		c.Next()
	}
}
