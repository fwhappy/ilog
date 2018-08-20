package ilog

import (
	"encoding/json"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/logs"
)

// LogConfig 日志配置
type LogConfig struct {
	AdapterConsole      bool     `toml:"adapter_console"`        // 是否输出到控制台
	ConsoleLevel        int      `toml:"console_level"`          // 输出到命令行的日志级别
	File                string   `toml:"file"`                   // 文件日志输出路径
	FileLevel           int      `toml:"file_level"`             // 文件日志级别
	EnableFuncCallDepth bool     `toml:"enable_func_call_depth"` // 是否输出行号和文件名
	Async               bool     `toml:"async"`                  // 是否异步输出,提升性能
	ChanLength          int      `toml:"chan_length"`            // 异步channel大小
	Rotate              bool     `toml:"rotate"`                 // 是否开启日志rotate
	Maxlines            int      `toml:"maxlines"`               // 单个文件行数限制
	Maxsize             int      `toml:"maxsize"`                // 单个文件大小限制
	Daily               bool     `toml:"daily"`                  // 是否按照每天rotat
	Maxdays             int      `toml:"maxdays"`                // 文件最多保留天数
	Multifile           bool     `toml:"multifile"`              // 是否分文件存储日志
	Separate            []string `toml:"separate"`               //分文件存储的日志类型
}

// CreateLogger 创建一个logger对象
func CreateLogger(cfgFile string) *BeeLogger {
	content, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}

	var cfg LogConfig
	if _, err := toml.Decode(string(content), &cfg); err != nil {
		panic(err)
	}

	log := NewLogger()

	// 设置异步输出
	if cfg.Async {
		log.Async(int64(cfg.ChanLength))
	}
	// 设置输出文件名、文件行数
	if cfg.EnableFuncCallDepth {
		log.EnableFuncCallDepth(true)
	}
	// 设置控制台输出
	if cfg.AdapterConsole {
		consoleConfig := make(map[string]int)
		consoleConfig["level"] = cfg.ConsoleLevel
		byt, _ := json.Marshal(consoleConfig)
		log.SetLogger(AdapterConsole, string(byt))
	}

	fileConfig := make(map[string]interface{})
	fileConfig["filename"] = cfg.File
	fileConfig["maxlines"] = cfg.Maxlines
	fileConfig["maxsize"] = cfg.Maxsize
	fileConfig["daily"] = cfg.Daily
	fileConfig["maxdays"] = cfg.Maxdays
	fileConfig["rotate"] = cfg.Rotate
	if cfg.Multifile {
		fileConfig["separate"] = cfg.Separate
		byt, _ := json.Marshal(fileConfig)
		log.SetLogger(logs.AdapterMultiFile, string(byt))
	} else {
		byt, _ := json.Marshal(fileConfig)
		log.SetLogger(AdapterFile, string(byt))
	}

	// 据说不这样做，会有一些性能问题
	log.SetLevel(cfg.FileLevel)

	return log
}
