/*
 * @PackageName: easy_dfs
 * @FileName: main.go
 * @Description: 项目入口
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:13:35
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:13:35
 */

package main

import (
	"easy_dfs/bootstrap"
	btsConfig "easy_dfs/config"
	"log"
)

func init() {
	// 初始化配置信息
	btsConfig.Initialize()
}

func main() {
	// 设置日志格式: 日期 时间 状态:成功/失败(成功则是绿色,失败则是红色)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 程序启动时打印
	log.Println("Server starting...")
	bootstrap.SetupLogger()
	bootstrap.SetupConfigDir()
	bootstrap.SetupRoute()
}
