// Package config 负责配置信息
package config

import (
	"fmt"
	"os"

	"github.com/spf13/cast"
	viperlib "github.com/spf13/viper"

	"github.com/boyane126/go-common/helpers"
)

// viper 库实例
var viper *viperlib.Viper

// ConfigFunc 动态加载配置信息
type ConfigFunc func() map[string]interface{}

var ConfigFuncs map[string]ConfigFunc

func init() {
	// 初始化 viper
	viper = viperlib.New()

	// 配置类型
	viper.SetConfigType("env")

	// 环境变量配置文件查找路径，相对于 main.go
	viper.AddConfigPath(".")

	// 设置环境变量前缀,用以区分 go 的系统环境变量
	viper.SetEnvPrefix("appenv")

	// 读取环境变量
	viper.AutomaticEnv()

	ConfigFuncs = make(map[string]ConfigFunc)
}

func InitConfig(env string) {
	// 加载 .env 环境变量
	loadEnv(env)

	// 注册配置信息
	loadConfig()
}

func loadEnv(env string) {
	// 默认加载 .env 文件，如果有传参 --env=name 的情况，加载 .env.name 文件
	envPath := ".env"
	if len(env) > 0 {
		filepath := ".env." + env
		if _, err := os.Stat(filepath); err == nil {
			// .env.testing 或 .env.stage
			envPath = filepath
		}
	}

	// 检查文件是否存在
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		fmt.Println("[Config] No .env file found, skipping loading config from env file")
		return
	}

	// 加载 env
	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.WatchConfig()
}

func loadConfig() {
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

// 读取环境变量，支持默认值
func Env(envName string, defaultVal ...interface{}) interface{} {
	if len(defaultVal) > 0 {
		return internalGet(envName, defaultVal[0])
	}
	return internalGet(envName)
}

func internalGet(path string, defaultVal ...interface{}) interface{} {
	// config 或者环境变量不存在的情况
	if !viper.IsSet(path) || helpers.Empty(viper.Get(path)) {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return nil
	}
	return viper.Get(path)
}

// Add 新增配置项
func Add(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

func Get(path string, defaultVal ...interface{}) string {
	return GetString(path, defaultVal...)
}

func GetString(path string, defaultVal ...interface{}) string {
	return cast.ToString(internalGet(path, defaultVal...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// GetFloat64 获取 float64 类型的配置信息
func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString 获取结构数据
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}

// 更新配置文件参数 例如：将MAIN_CTRL_ADDR=localhost:50051 更新为MAIN_CTRL_ADDR=127.0.0.1:50051
func UpdateConfig(path string, value interface{}) {
	viper.Set(path, value)
	viper.WriteConfig()
}
