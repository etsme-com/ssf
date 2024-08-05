package persistence

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/etsme-com/ssf/v2/base/config"
	"github.com/etsme-com/ssf/v2/define"
	"github.com/etsme-com/ssf/v2/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var persistenceLogger = logger.Logger.WithFields(logrus.Fields{"module": "persistence"})

// //////////////// data persistence //////////////////
// 持久化模式
var (
	NeverLostPersistence        = 0 // 永不丢失，		EMMC
	FactoryResetLostPersistence = 1 // 恢复出厂设置丢失，SSD
	UnboundLostPersistence      = 2 // 解绑丢失，		SSD
	PoweroffLostPersistence     = 3 // 掉电丢失，		保留内存
	RebootLostPersistence       = 4 // 重启丢失，		tmp
)

var (
	PersistenceCommonPhdDevState    = "PhdDevState"
	PersistenceCommonBound          = "BoundState"
	PersistenceCommonPhdClusterInfo = "PHDClusterInfo"
	PersistenceCommonModlesReady    = "AIModlesReady"
	PersistenceCommoUpgradeRunning  = "UpgradeRunning"
)

type persistenceConfigVer struct {
	Version int `yaml:"version"`
}

type persistenceConfigV1 struct {
	Version  int               `yaml:"version"`
	Elements map[string][]byte `yaml:"element"`
}

type persistenceConfig struct {
	Version  int               `yaml:"version"`
	Elements map[string]string `yaml:"element"`
}

var commonStr = "common"
var commonLock = "PERSISTENCE:COMMON"

var lock sync.Mutex

var (
	persistenceVersionV1 = 100
	persistenceVersion   = 300
)

const (
	fileNew = ".new"
	fileMd5 = ".md5"
)

var persistenceProtectPath string

func init() {
	if config.SSFConfig.Platform == define.PlatformCloud {
		persistenceLogger.Infoln("cloud not support.")
		return
	}

	persistenceProtectPath = fmt.Sprintf("%s/persistence/protect", config.SSFConfig.Storage.SysCfgPath)
	persistenceLogger.Infof("persistenceProtectPath=%s", persistenceProtectPath)

	err := os.MkdirAll(persistenceProtectPath, os.ModePerm)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":  err,
			"path": persistenceProtectPath,
		}).Errorln("MkdirAll fail.")
	}
}

func CleanProtectDir() {
	persistenceProtectPath = fmt.Sprintf("%s/persistence/protect/*", config.SSFConfig.Storage.SysCfgPath)
	cmd := fmt.Sprintf("rm %s -rf", persistenceProtectPath)
	persistenceLogger.Infoln(cmd)
	output, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err": err,
		}).Errorln("rm failed:" + string(output))
	}
}

func StoreServiceCustomPersistenceData(name string, mode int, key string, value interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	filename := fmt.Sprintf("%s-%s", config.ServiceConfig.Name, name)
	return storePersistenceData(filename, mode, key, value)
}

func LoadServiceCustomPersistenceAllData(name string, mode int) (map[string]string, error) {
	lock.Lock()
	defer lock.Unlock()

	filename := fmt.Sprintf("%s-%s", config.ServiceConfig.Name, name)
	return loadPersistenceAllData(filename, mode)
}

func LoadServiceCustomPersistenceData(name string, mode int, key string) ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()

	filename := fmt.Sprintf("%s-%s", config.ServiceConfig.Name, name)
	return loadPersistenceData(filename, mode, key)
}

func DelServiceCustomPersistenceData(name string, mode int, key string) error {
	lock.Lock()
	defer lock.Unlock()

	filename := fmt.Sprintf("%s-%s", config.ServiceConfig.Name, name)
	return delPersistenceData(filename, mode, key)
}

func StoreServicePersistenceData(mode int, key string, value interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	return storePersistenceData(config.ServiceConfig.Name, mode, key, value)
}

func LoadServicePersistenceData(mode int, key string) ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()

	return loadPersistenceData(config.ServiceConfig.Name, mode, key)
}

func DelServicePersistenceData(mode int, key string) error {
	lock.Lock()
	defer lock.Unlock()

	return delPersistenceData(config.ServiceConfig.Name, mode, key)
}

// 无锁load
func LoadCommonPersistenceDataNoLock(mode int, key string) ([]byte, error) {
	var err error

	data, err := loadPersistenceData(commonStr, mode, key)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err": err,
		}).Errorln("loadPersistenceData fail.")
	}

	return data, err
}
func getDataPath(mode int) (string, error) {
	var datapath string

	switch mode {
	case NeverLostPersistence:
		datapath = fmt.Sprintf("%s/persistence", config.SSFConfig.Storage.SysCfgPath)
	case FactoryResetLostPersistence:
		datapath = fmt.Sprintf("%s/persistence", config.SSFConfig.Storage.MetadataPath)
	case UnboundLostPersistence:
		datapath = fmt.Sprintf("%s/persistence", config.SSFConfig.Storage.MetadataPath)
	case RebootLostPersistence:
		datapath = fmt.Sprintf("%s/persistence", config.SSFConfig.Storage.TmpPath)
	default:
		datapath = fmt.Sprintf("%s/persistence", config.SSFConfig.Storage.TmpPath)
	}

	err := os.MkdirAll(datapath, os.ModePerm)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":  err,
			"path": datapath,
		}).Errorln("MkdirAll data failed.")
	}

	persistenceProtectPath = fmt.Sprintf("%s/persistence/protect", config.SSFConfig.Storage.SysCfgPath)
	protectPath := fmt.Sprintf("%s%s", persistenceProtectPath, datapath)
	err = os.MkdirAll(protectPath, os.ModePerm)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":                    err,
			"persistenceProtectPath": persistenceProtectPath,
			"path":                   protectPath,
		}).Errorln("MkdirAll Pprotect dir failed.")
	}

	return datapath, nil
}

func getNewFilePath(path string) string {
	return fmt.Sprintf("%s%s%s", persistenceProtectPath, path, fileNew)
}

func getMd5FilePath(path string) string {
	return fmt.Sprintf("%s%s%s", persistenceProtectPath, path, fileMd5)
}

func readFile(path string) ([]byte, error) {
	// 1. 读取 真实文件内容
	data, err := ioutil.ReadFile(path)

	// 掉电丢失配置，无需保护
	if strings.Contains(path, "/tmp/persistence") {
		return data, err
	}

	// 2. 读取 md5 文件
	md5FilePath := getMd5FilePath(path)
	fileMd5Data, errmd5 := ioutil.ReadFile(md5FilePath)

	// 如果 md5文件打开失败，则不做对比操作
	if errmd5 == nil {
		// 3. 计算真实文件的MD5
		md5Data := md5.Sum(data)

		// 4. 真实文件的md5 和 .md5文件内容不一致，取 .new 文件内容
		if bytes.Compare(fileMd5Data, md5Data[:]) != 0 {
			persistenceLogger.WithFields(logrus.Fields{
				"md5Data":     md5Data,
				"fileMd5Data": fileMd5Data,
			}).Infoln("md5 mismatch, read .new file")

			newFilePath := getNewFilePath(path)
			newData, err := ioutil.ReadFile(newFilePath)
			if err == nil {
				newMd5Data := md5.Sum(newData)
				// 5. .new 文件的md5 和 .md5文件内容不一致，修正md5文件
				if bytes.Compare(fileMd5Data, newMd5Data[:]) != 0 {
					persistenceLogger.WithFields(logrus.Fields{
						"path": md5FilePath,
					}).Infoln("fix md5.")

					f, err := os.OpenFile(md5FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0777)
					if err != nil {
						persistenceLogger.WithFields(logrus.Fields{
							"err":  err,
							"file": md5FilePath,
						}).Errorln("OpenFile fail.")
					}

					_, err = f.Write(newMd5Data[:])
					if err != nil {
						persistenceLogger.WithFields(logrus.Fields{
							"err": err,
						}).Errorln("Write fail.")
					}

					f.Close()
				}

				// 6. 覆盖写入真实文件
				persistenceLogger.WithFields(logrus.Fields{
					"path": newFilePath,
				}).Infoln("fix file.")

				f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0777)
				if err != nil {
					persistenceLogger.WithFields(logrus.Fields{
						"err":  err,
						"file": newFilePath,
					}).Errorln("OpenFile fail.")
				}

				_, err = f.Write(newData)
				if err != nil {
					persistenceLogger.WithFields(logrus.Fields{
						"err": err,
					}).Errorln("Write fail.")
				}

				f.Close()

				// 6. 返回 .new 文件内容
				return newData, nil
			}
		}
	}

	// 如果md5 文件不存在，则直接返回真实文件内容
	return data, err
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	// 掉电丢失配置，无需保护
	if !strings.Contains(path, "/tmp/persistence") {
		// 1. 新内容先写入 .new 文件
		newFilePath := getNewFilePath(path)
		newFile, err := os.OpenFile(newFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, perm)
		if err != nil {
			persistenceLogger.WithFields(logrus.Fields{
				"err":  err,
				"file": newFilePath,
			}).Errorln("OpenFile fail.")
			return err
		}

		_, err = newFile.Write(data)
		if err != nil {
			persistenceLogger.WithFields(logrus.Fields{
				"err": err,
			}).Errorln("Write fail.")
			return err
		}

		err = newFile.Close()
		if err != nil {
			persistenceLogger.WithFields(logrus.Fields{
				"err":  err,
				"file": newFilePath,
			}).Errorln("Close fail.")
			return err
		}

		// 2.计算MD5写入 .md5 文件
		md5FilePath := getMd5FilePath(path)
		md5File, err := os.OpenFile(md5FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, perm)
		if err != nil {
			persistenceLogger.WithFields(logrus.Fields{
				"err":  err,
				"file": md5FilePath,
			}).Errorln("OpenFile fail.")
			return err
		}

		md5Data := md5.Sum(data)
		_, err = md5File.Write(md5Data[:])
		if err != nil {
			persistenceLogger.WithFields(logrus.Fields{
				"err": err,
			}).Errorln("Write fail.")
			return err
		}

		err = md5File.Close()
		if err != nil {
			persistenceLogger.WithFields(logrus.Fields{
				"err":  err,
				"file": md5FilePath,
			}).Errorln("Close fail.")
			return err
		}
	}

	// 3.写入真实文件
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, perm)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":  err,
			"file": path,
		}).Errorln("OpenFile fail.")
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err": err,
		}).Errorln("Write fail.")
	}

	err = f.Close()
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":  err,
			"file": path,
		}).Errorln("Close fail.")
		return err
	}

	return err
}

func repair(filePath string) {
	yamlfile, err := readFile(filePath)
	if err == nil {
		var pcfgver persistenceConfigVer
		err = yaml.Unmarshal(yamlfile, &pcfgver)
		if err != nil {
			persistenceLogger.WithFields(logrus.Fields{
				"err":      err,
				"filePath": filePath,
			}).Errorln("Unmarshal yaml fail.")
			return
		}

		// 如果 是 V1 版本，做数据转换
		if pcfgver.Version == persistenceVersionV1 {
			var pcfgv1 persistenceConfigV1
			err = yaml.Unmarshal(yamlfile, &pcfgv1)
			if err != nil {
				persistenceLogger.WithFields(logrus.Fields{
					"err":      err,
					"filePath": filePath,
				}).Errorln("Unmarshal yaml fail.")
				return
			}

			pcfg := persistenceConfig{
				Version:  persistenceVersion,
				Elements: make(map[string]string, len(pcfgv1.Elements)),
			}

			for i, v := range pcfgv1.Elements {
				pcfg.Elements[i] = string(v)
			}

			pcfgdata, err := yaml.Marshal(&pcfg)
			if err != nil {
				persistenceLogger.WithFields(logrus.Fields{
					"error":    err,
					"filePath": filePath,
				}).Errorln("Marshal yaml failed.")
				return
			}

			err = writeFile(filePath, pcfgdata, 0777)
			if err != nil {
				persistenceLogger.WithFields(logrus.Fields{
					"error":    err,
					"filePath": filePath,
				}).Errorln("writeFile failed.")
			}
		}
	}
}

func storePersistenceData(class string, mode int, key string, value interface{}) error {
	var pcfg persistenceConfig

	datapath, err := getDataPath(mode)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err": err,
		}).Errorln("getDataPath fail.")
		return err
	}

	filePath := fmt.Sprintf("%s/%s.yaml", datapath, class)

	// 数据格式修正
	repair(filePath)

	yamlfile, err := readFile(filePath)
	if err == nil {
		err = yaml.Unmarshal(yamlfile, &pcfg)
		if err != nil {
			persistenceLogger.WithFields(logrus.Fields{
				"err":      err,
				"filePath": filePath,
				"key":      key,
			}).Errorln("Unmarshal yaml fail, discard data, repair below")

			// 如果读取配置文件后，解析失败，说明文件已被损坏，不返回错误，后面修复
			// return err
		}
	}

	if pcfg.Elements == nil {
		persistenceLogger.Infoln("File Corrupted, ", filePath)

		pcfg = persistenceConfig{
			Version:  persistenceVersion,
			Elements: make(map[string]string),
		}
	}

	data, err := json.Marshal(value)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"error":    err,
			"filePath": filePath,
			"key":      key,
		}).Errorln("Marshal json failed.")
		return err
	}

	pcfg.Elements[key] = string(data)
	pcfgdata, err := yaml.Marshal(&pcfg)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"error":    err,
			"filePath": filePath,
			"key":      key,
		}).Errorln("Marshal yaml failed.")
		return err
	}

	return writeFile(filePath, pcfgdata, 0777)
}

func loadPersistenceAllData(class string, mode int) (map[string]string, error) {
	datapath, err := getDataPath(mode)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err": err,
		}).Errorln("getDataPath fail.")
		return nil, err
	}

	filePath := fmt.Sprintf("%s/%s.yaml", datapath, class)

	// 数据格式修正
	repair(filePath)

	yamlfile, err := readFile(filePath)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":      err,
			"filePath": filePath,
		}).Errorln("ReadFile fail.")

		if os.IsNotExist(err) {
			return nil, errors.New("Key Not Exist")
		}

		return nil, err
	}

	var pcfg persistenceConfig
	err = yaml.Unmarshal(yamlfile, &pcfg)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":      err,
			"filePath": filePath,
		}).Errorln("Unmarshal fail.")
		return nil, err
	}

	if pcfg.Elements == nil {
		persistenceLogger.Infoln("File Corrupted, ", filePath)

		pcfg = persistenceConfig{
			Version:  persistenceVersion,
			Elements: make(map[string]string),
		}
	}

	return pcfg.Elements, nil
}

func loadPersistenceData(class string, mode int, key string) ([]byte, error) {
	elements, err := loadPersistenceAllData(class, mode)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err": err,
			"key": key,
		}).Errorln("loadPersistenceAllData fail.")
		return nil, err
	}

	if value, ok := elements[key]; ok {
		return []byte(value), nil
	}

	return nil, errors.New("Key Not Exist")
}

func delPersistenceData(class string, mode int, key string) error {
	var pcfg persistenceConfig

	datapath, err := getDataPath(mode)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err": err,
		}).Errorln("getDataPath fail.")
		return err
	}

	filePath := fmt.Sprintf("%s/%s.yaml", datapath, class)

	// 数据格式修正
	repair(filePath)

	yamlfile, err := readFile(filePath)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":      err,
			"filePath": filePath,
			"key":      key,
		}).Errorln("ReadFile fail.")

		if os.IsNotExist(err) {
			return errors.New("Key Not Exist")
		}

		return err
	}

	err = yaml.Unmarshal(yamlfile, &pcfg)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"err":      err,
			"filePath": filePath,
			"key":      key,
		}).Errorln("Unmarshal yaml fail.")
		return err
	}

	if pcfg.Elements == nil {
		persistenceLogger.Infoln("File Corrupted, ", filePath)

		pcfg = persistenceConfig{
			Version:  persistenceVersion,
			Elements: make(map[string]string),
		}
	}

	if _, ok := pcfg.Elements[key]; !ok {
		return errors.New("Key Not Exist")
	}

	delete(pcfg.Elements, key)
	pcfgdata, err := yaml.Marshal(&pcfg)
	if err != nil {
		persistenceLogger.WithFields(logrus.Fields{
			"error":    err,
			"filePath": filePath,
			"key":      key,
		}).Errorln("Marshal yaml failed.")
		return err
	}

	return writeFile(filePath, pcfgdata, 0777)
}
