/*
 * @PackageName: bootstrap
 * @FileName: config_dir.go
 * @Description: 配置文件
 * @Author: gabbymrh
 * @Date: 2024-07-18 11:45:17
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 11:45:17
 */

package bootstrap

import (
	"easy_dfs/pkg/config"
	"os"
)

// 引导生成配置目录及文件
func SetupConfigDir() {
	// 配置目录
	configDir := "config"
	// 文件存储目录
	storageDir := "storage"

	appEnv := config.Get("app.env")

	if appEnv != "prod" {
		configDir = "tmp/config"
		storageDir = "tmp/storage"
	}
	// access key配置文件
	accessKeyConfig := configDir + "/access_key.json"
	// bucket配置文件
	bucketConfig := configDir + "/bucket.json"

	// 创建配置目录
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		panic(err)
	}
	// 创建文件存储目录
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		panic(err)
	}
	// 创建access key配置文件并写入空数组
	if _, err := os.Stat(accessKeyConfig); os.IsNotExist(err) {
		file, err := os.Create(accessKeyConfig)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 写入空数组数据
		_, err = file.WriteString("[]\n")
		if err != nil {
			panic(err)
		}
	}
	// 创建bucket配置文件并写入空数组
	if _, err := os.Stat(bucketConfig); os.IsNotExist(err) {
		file, err := os.Create(bucketConfig)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 写入空数组数据
		_, err = file.WriteString("[]\n")
		if err != nil {
			panic(err)
		}
	}
}
