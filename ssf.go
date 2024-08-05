package ssf

import (
	"gitlab.etsme.com/cmcc-stable-lc/osware/ssf/ssf/v2/base/config"
	"gitlab.etsme.com/cmcc-stable-lc/osware/ssf/ssf/v2/logger"
)

// SSFVersion 微服务框架版本
const SSFVersion = "v2.0.4"

// SSFConfig 公共相关配置
var SSFConfig = &config.SSFConfig

// Logger  日志操作全局对象引用
var Logger = logger.Logger

// init 初始化函数
func init() {
	Logger.Infof("ssf init success, version = %s.", SSFVersion)
}

func GetVersion() string {
	return SSFVersion
}
