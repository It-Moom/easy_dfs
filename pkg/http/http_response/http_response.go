/*
 * @PackageName: http_response
 * @FileName: http_response.go
 * @Description: HTTP响应
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:52:26
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:52:26
 */

package http_response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 返回数据结构体
type ResponseData struct {
	Code    string      `json:"code"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
}

// 返回体封装
func Response(ctx *gin.Context, code string, success bool, message string, data interface{}, errors error) {
	myErrors := make([]string, 0)
	if errors != nil {
		myErrors = append(myErrors, errors.Error())
	}
	ctx.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Success: success,
		Data:    data,
		Message: message,
		Errors:  myErrors,
	})
	// 终止请求处理
	ctx.Abort()
}
