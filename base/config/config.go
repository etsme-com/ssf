package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type ssfYamlConfig struct {
	Platform int64 `yaml:"platform"`

	EtsmePath string `yaml:"etsmePath"`

	Storage struct {
		SysCfgPath    string              `yaml:"sysCfgPath"`
		SysBackupPath string              `yaml:"sysBackupPath"`
		TmpPath       string              `yaml:"tmpPath"`
		BackupPath    string              `yaml:"backupPath"`
		LogPath       string              `yaml:"logPath"`
		MetadataPath  string              `yaml:"metadataPath"`
		CachePath     string              `yaml:"cachePath"`
		BlockDevPath  []string            `yaml:"blockDevPath"`
		Disks         []map[string]string `yaml:"disks"`
	}

	Route []map[string]string `yaml:"route"`

	Cloud struct {
		Services []map[string]string `yaml:"services"`
	}
}

type serviceYamlConfig struct {
	Name    string
	Version string

	Logger struct {
		LogPath          string `yaml:"logPath"`
		LogLevel         string `yaml:"logLevel"`
		LogOutToFile     bool   `yaml:"logOutToFile"`
		LogRotationSize  int64  `yaml:"logRotationSize"`
		LogRotationCount uint   `yaml:"logRotationCount"`
		LogEnGID         bool   `yaml:"logEnGID"`
		LogHasCaller     bool   `yaml:"logHasCaller"`
	}
}

// SSFConfig 微服务通用配置信息
var SSFConfig ssfYamlConfig

// ServiceConfig 微服务特性配置信息
var ServiceConfig serviceYamlConfig

// init 获取yaml配置文件，结构化到SSFConfig和ServiceConfig中
func init() {
	ssfYamlFile, err := ioutil.ReadFile("/etc/ssf/SSFConfig.yaml")
	if err != nil {
		ssfYamlFile, err = ioutil.ReadFile("/opt/etsme/etc/ssf/SSFConfig.yaml")
		if err != nil {
			fmt.Printf("read /etc/ssf/SSFConfig.yaml and /opt/etsme/etc/ssf/SSFConfig.yaml failed, %s\n", err)
			return
		}
	}

	err = yaml.Unmarshal(ssfYamlFile, &SSFConfig)
	if err != nil {
		fmt.Printf("Unmarshal SSFConfig.yaml fail, %s\n", err)
		return
	}

	namel := strings.Split(os.Args[0], "/")
	serviceName := namel[len(namel)-1]

	ServiceConfig.Logger.LogLevel = "info"
	ServiceConfig.Logger.LogOutToFile = true
	ServiceConfig.Logger.LogRotationSize = 10485760
	ServiceConfig.Logger.LogRotationCount = 14
	ServiceConfig.Logger.LogEnGID = false
	ServiceConfig.Logger.LogHasCaller = true
	ServiceConfig.Name = serviceName
}

// 获取 产品平台
func GetPlatform() int64 {
	if SSFConfig.Platform < 100 {
		return SSFConfig.Platform
	}

	return (SSFConfig.Platform / 100) * 100
}
