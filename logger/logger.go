// Package logger 日志处理包
// Author :  陈焱    2020/08/01
package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	rotatelogs "github.com/etsme-com/file-rotatelogs"
	"github.com/etsme-com/ssf/base/config"
	log "github.com/sirupsen/logrus"
)

var logger = log.New()
var opLogger = log.New()
var chainLogger = log.New()
var evtLogger = log.New()
var alarmLogger = log.New()

// Logger 调试日志全局对象引用
var Logger = logger.WithFields(log.Fields{})

//var Logger = logger

// OpLogger 操作日志全局对象引用
var OpLogger = opLogger

// ChainLogger 调用链日志全局对象引用
var ChainLogger = chainLogger

// EvtLogger 事件日志全局对象引用
var EvtLogger = evtLogger

// AlarmLogger 用户告警日志全局对象引用
var AlarmLogger = alarmLogger

// ssfFormatter 自定义 formatter
type ssfFormatter struct {
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// Format implement the Formatter interface
func (mf *ssfFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	fieldsStr := ""
	if len(entry.Data) > 0 {
		mjson, _ := json.Marshal(entry.Data)
		fieldsStr = string(mjson)
	}

	if entry.HasCaller() {
		fileList := strings.Split(entry.Caller.File, "/")
		file := fileList[len(fileList)-1]

		funcList := strings.Split(entry.Caller.Function, "/")
		function := funcList[len(funcList)-1]

		if true {
			b.WriteString(fmt.Sprintf("[%s] [%d] [%s] [%s %d %s] %s %s\n",
				entry.Time.Format("2006-01-02 15:04:05.000"),
				getGID(),
				entry.Level,
				file,
				entry.Caller.Line,
				function,
				fieldsStr,
				entry.Message),
			)
		} else {
			b.WriteString(fmt.Sprintf("[%s] [%s] [%s %d %s] %s %s\n",
				entry.Time.Format("2006-01-02 15:04:05.000"),
				entry.Level,
				file,
				entry.Caller.Line,
				function,
				fieldsStr,
				entry.Message),
			)
		}
	} else {
		b.WriteString(fmt.Sprintf("[%s] [%s] %s %s\n",
			entry.Time.Format("2006-01-02 15:04:05.000"),
			entry.Level,
			fieldsStr,
			entry.Message),
		)
	}

	return b.Bytes(), nil
}

// ssfOpFormatter 自定义 formatter
type ssfOpFormatter struct {
}

// Format implement the Formatter interface
func (mf *ssfOpFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	fieldsStr := ""
	if len(entry.Data) > 0 {
		mjson, _ := json.Marshal(entry.Data)
		fieldsStr = string(mjson)
	}

	b.WriteString(fmt.Sprintf("[%s] %s %s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		fieldsStr,
		entry.Message),
	)
	return b.Bytes(), nil
}

// ssfAlarmFormatter 自定义 formatter
type ssfAlarmFormatter struct {
}

// Format implement the Formatter interface
func (mf *ssfAlarmFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(fmt.Sprintf("%s\n",
		entry.Message),
	)
	return b.Bytes(), nil
}

// init logrus 初始化
func init() {
	formatter := &ssfFormatter{}

	// 调试日志
	logger.SetFormatter(formatter)
	logger.SetReportCaller(true)

	// 设置默认输出
	SetOutput()

	SetLogLevel(config.ServiceConfig.Logger.LogLevel)
	Logger.Infoln("logger init success.")

}

// SetOutput 设置输出路径
func SetOutput() {
	serviceLogPath := fmt.Sprintf("%s/services", config.SSFConfig.Storage.LogPath)

	//fmt.Println(serviceLogPath)
	err := os.MkdirAll(serviceLogPath, 0777)
	if err != nil {
		fmt.Println("MkdirAll failed.")
	}

	SwitchOutputPath(serviceLogPath)
}

// SwitchOutputPath 切换输出路径
func SwitchOutputPath(path string) {
	logdir := fmt.Sprintf("%s/%s", path, config.ServiceConfig.Name)
	file := fmt.Sprintf("%s/%s.log", logdir, config.ServiceConfig.Name)

	// 日志轮转相关函数
	// WithLinkName 为最新的日志建立软连接
	// WithRotationTime 设置日志分割的时间，隔多久分割一次, 默认24小时更新一个文件
	// WithRotationSize 设置日志的分割的容量
	// WithMaxAge 和 WithRotationCount 二者只能设置一个
	// WithMaxAge 设置文件清理前的最长保存时间
	// WithRotationCount 设置文件清理前最多保存的个数
	writer, err := rotatelogs.New(
		file+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(file),
		rotatelogs.WithRotationCount(config.ServiceConfig.Logger.LogRotationCount),
		rotatelogs.WithRotationSize(config.ServiceConfig.Logger.LogRotationSize),
	)
	if err == nil {
		logger.SetOutput(writer)
	} else {
		fmt.Println("writer init fail.")
		logger.SetOutput(os.Stdout)
	}
}

// SetLogLevel 配置log级别
func SetLogLevel(level string) {
	switch level {
	case "trace":
		logger.SetLevel(log.TraceLevel)
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warning":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "fatal":
		logger.SetLevel(log.FatalLevel)
	case "panic":
		logger.SetLevel(log.PanicLevel)
	default:
		logger.SetLevel(log.WarnLevel)
	}
}
