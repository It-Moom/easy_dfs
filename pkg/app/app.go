/*
 * @PackageName: app
 * @FileName: app.go
 * @Description: 应用文件
 * @Author: gabbymrh
 * @Date: 2024-07-18 14:21:16
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 14:21:16
 */

package app

import (
	"easy_dfs/pkg/config"
	"time"
)

// IsLocal 是否本地环境
func IsLocal() bool {
	return config.Get("app.env") == "local"
}

// IsProduction 是否生产环境
func IsProduction() bool {
	return config.Get("app.env") == "prod"
}

// IsTesting 是否测试环境
func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(chinaTimezone)
}

// URL 传参 path 拼接站点的 URL
func URL(path string) string {
	return config.Get("app.url") + path
}
