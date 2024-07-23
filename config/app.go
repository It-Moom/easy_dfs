/*
 * @PackageName: config
 * @FileName: app.go
 * @Description: 应用配置文件
 * @Author: gabbymrh
 * @Date: 2024-07-18 12:20:16
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 12:20:16
 */

package config

import (
	"easy_dfs/pkg/config"
)

func init() {
	config.Add("app", func() map[string]interface{} {
		return map[string]interface{}{
			// 应用名称
			"name": config.Env("app.name", "EasyDFS"),
			// 当前环境，用以区分多环境
			"env": config.Env("app.env", "prod"),
			// 是否进入调试模式
			"debug": config.Env("app.debug", false),
			// App url
			"url": config.Env("app.url", "http://localhost"),
			// 应用服务端口
			"port": config.Env("app.port", "18088"),
			// 设置时区，JWT 里会使用，日志记录里也会使用到
			"timezone": config.Env("app.timezone", "Asia/Shanghai"),
			// API 域名，未设置的话所有 API URL 加 api 前缀，如 http://domain.com/api/v1/users
			"api_domain": config.Env("app.api-url", ""),
		}
	})
}
