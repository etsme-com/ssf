package envcfg

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"github.com/etsme-com/ssf/base/config"
	"github.com/etsme-com/ssf/define"
	"github.com/etsme-com/ssf/logger"
	"github.com/etsme-com/ssf/persistence"
	"github.com/sirupsen/logrus"
)

// Cenvcfg 保存系统环境配置
type Cenvcfg struct {
	NetworkEnv int // 网络环境
}

// Envcfg 全局描述符
var Envcfg Cenvcfg
var envFilePath = config.SSFConfig.Storage.SysCfgPath + "/env"

var envLogger = logger.Logger.WithFields(logrus.Fields{"module": "env"})

func init() {
	Envcfg = Cenvcfg{}
}

var PersistenceCommonPhEnv = "PhdEnv"

func readEnvFromPersis() (string, error) {
	envStr := ""

	data, err := persistence.LoadCommonPersistenceDataNoLock(persistence.UnboundLostPersistence, PersistenceCommonPhEnv)
	if err != nil {
		envLogger.WithFields(logrus.Fields{
			"err":                    err,
			"PersistenceCommonPhEnv": PersistenceCommonPhEnv,
		}).Infoln("LoadCommonPersistenceData failed.")
	} else {
		err = json.Unmarshal(data, &envStr)
		if err != nil {
			envLogger.WithFields(logrus.Fields{
				"err": err,
			}).Errorln("Unmarshal failed.")
		}
	}

	return envStr, err
}

func GetRunEnv() (int, string) {
	// 尝试从持久化配置里获取
	envStr, err := readEnvFromPersis()
	if err != nil {
		envLogger.WithFields(logrus.Fields{
			"err": err,
		}).Infoln("readEnvFromPersis failed.")
		return -1, ""
	}

	var env int
	if strings.Contains(string(envStr), "testing2") {
		env = define.Testing2RunEnv
	} else if strings.Contains(string(envStr), "testing") {
		env = define.TestingRunEnv
	} else if strings.Contains(string(envStr), "development") {
		env = define.DevelopmentRunEnv
	} else {
		env = define.ProductionRunEnv
		envStr = "product"
	}

	return env, envStr
}

func GetRunEnvIntFromStr(envStr string) int {
	var env int
	if strings.Contains(string(envStr), "testing2") {
		env = define.Testing2RunEnv
	} else if strings.Contains(string(envStr), "testing") {
		env = define.TestingRunEnv
	} else if strings.Contains(string(envStr), "development") {
		env = define.DevelopmentRunEnv
	} else {
		// 默认都是生产
		env = define.ProductionRunEnv
	}

	return env
}

func GetCloudDomainName(env string, key string) ([]string, error) {
	var envStr string

	if len(env) > 0 {
		envStr = env
	} else {
		_, envStr = GetRunEnv()
	}

	for _, service := range config.SSFConfig.Cloud.Services {
		name, ok := service["name"]
		if ok && strings.Compare(name, key) == 0 {
			if host, ok := service[envStr]; ok {
				parts := strings.Split(host, "&&")
				return parts, nil
			}
		}
	}

	return nil, errors.New("Err Parameter")
}

func strIn(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	//index的取值：[0,len(str_array)]
	//需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
