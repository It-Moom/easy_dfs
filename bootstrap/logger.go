/*
 * @PackageName: bootstrap
 * @FileName: logger.go
 * @Description: 日志配置
 * @Author: gabbymrh
 * @Date: 2024-07-18 14:38:39
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 14:38:39
 */

package bootstrap

import (
	"easy_dfs/pkg/config"
	"easy_dfs/pkg/logger"
)

// SetupLogger 初始化 Logger
func SetupLogger() {
	logger.InitLogger(
		config.GetString("log.filename", "logs/logs.log"),
		config.GetInt("log.max_size", 64),
		config.GetInt("log.max_backup", 5),
		config.GetInt("log.max_age", 30),
		config.GetBool("log.compress", false),
		config.GetString("log.type", "daily"),
		config.GetString("log.level", "error"),
	)
}
