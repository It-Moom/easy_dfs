/*
 * @PackageName: response_code
 * @FileName: response_code.go
 * @Description: 响应码枚举
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:51:23
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:51:23
 */

package response_code

const (
	// 参数有误
	PARAM_ERROR = "10001"
	// 操作成功
	REQUEST_SUCCESS = "20000"
	// 拒绝访问
	REQUEST_DENIED = "40001"
	// 查询为空
	QUERY_EMPTY = "40004"
	// Token无效
	TOKEN_INVALID = "40005"
	// 请求频繁
	REQUEST_FREQUENT = "40029"
	// 操作失败
	REQUEST_FAILS = "50000"
)
