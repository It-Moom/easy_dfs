/*
 * @PackageName: config
 * @FileName: config_dir.go
 * @Description: 配置文件
 * @Author: gabbymrh
 * @Date: 2024-07-18 12:16:38
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 12:16:38
 */

package config

import (
	"easy_dfs/pkg/helpers"
	"fmt"
	"github.com/spf13/cast"
	viperlib "github.com/spf13/viper"
	"os"
)

// Viper Viper 库实例
var viper *viperlib.Viper

// ConfigFunc 动态加载配置信息
type ConfigFunc func() map[string]interface{}

// ConfigFuncs 先加载到此数组，loadConfig 再动态生成配置信息
var ConfigFuncs map[string]ConfigFunc

func init() {
	// 1. 初始化 Viper 库
	viper = viperlib.New()
	// 2. 设置文件名称
	viper.SetConfigName("app.yml")
	// 3. 配置类型，支持 "json", "toml", "yaml", "yml", "properties",
	//             "props", "prop", "env", "dotenv"
	viper.SetConfigType("yml")
	// 4. 环境变量配置文件查找的路径，相对于 main.go
	viper.AddConfigPath(".")

	// 5. 开始读根目录下的 app.yml 文件，读不到会报错
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 6. 设置环境变量前缀，用以区分 Go 的系统环境变量
	viper.SetEnvPrefix("appenv")
	// 7. Viper.Get() 时，优先读取环境变量
	viper.AutomaticEnv()

	ConfigFuncs = make(map[string]ConfigFunc)
}

// InitConfig 初始化配置信息，完成对环境变量以及 config 信息的加载
func InitConfig(env string) {
	// 1. 加载环境变量
	loadEnv(env)
	// 2. 注册配置信息
	loadConfig()
}

// loadConfig 加载配置
func loadConfig() {
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

// loadEnv 加载配置文件
func loadEnv(envSuffix string) {

	// 默认加载 app.yml 文件，如果有传参 --env=name 的话，加载 app.yml.name 文件
	envPath := "app.yml"
	if len(envSuffix) > 0 {
		filepath := "app.yml." + envSuffix
		if _, err := os.Stat(filepath); err == nil {
			// 如 app.yml.testing 或 app.yml.stage
			envPath = filepath
		}
	}

	// 加载 yml
	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 监控 app.yml 文件，变更时重新加载
	viper.WatchConfig()
}

// Env 读取环境变量，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return internalGet(envName, defaultValue[0])
	}
	return internalGet(envName)
}

// Add 新增配置项
func Add(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

// ViperUnMarshal Viper 反序列化
func ViperUnMarshal(v interface{}) interface{} {
	err := viper.Unmarshal(v)
	if err != nil {
		fmt.Println("Unmarshal Error:", err)
		return nil
	}
	return v
}

// Get 获取配置项
// 第一个参数 path 允许使用点式获取，如：app.name
// 第二个参数允许传参默认值
func Get(path string, defaultValue ...interface{}) string {
	return GetString(path, defaultValue...)
}

// internalGet 获取配置信息
func internalGet(path string, defaultValue ...interface{}) interface{} {
	// config 或者环境变量不存在的情况
	if !viper.IsSet(path) || helpers.Empty(viper.Get(path)) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(path)
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetFloat64 获取 float64 类型的配置信息
func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// ViperGetStringMap 获取结构数据
func ViperGetStringMap(path string) map[string]interface{} {
	return viper.GetStringMap(path)
}

// GetStringMapString 获取结构数据(字符串)
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}
