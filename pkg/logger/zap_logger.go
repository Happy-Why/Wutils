package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	*zap.Logger
}

var initLog = new(InitStruct)

type InitStruct struct {
	LogSavePath   string // 保存路径
	LogFileExt    string // 日志文件后缀
	MaxSize       int    // 备份的大小(M)
	MaxBackups    int    // 最大备份数
	MaxAge        int    // 最大备份天数
	Compress      bool   // 是否压缩过期日志
	LowLevelFile  string // 低级别文件名
	HighLevelFile string // 高级别文件名
}

func NewLogger(x *InitStruct, level string) *Log {
	initLog = x
	//判断日志级别是否大于或等于 ErrorLevel
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	//判断日志级别是否小于 ErrorLevel 且大于或等于 DebugLevel, info和debug级别,debug级别是最低的
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})
	// 多个日志文件
	var cores []zapcore.Core
	// 调用 getLogWriter 函数来获取低级别和高级别日志的写入器。这些写入器用于将日志写入文件。
	lowFileWriteSyncer := getLogWriter(initLog.LogSavePath + initLog.LowLevelFile + initLog.LogFileExt)
	highFileWriteSyncer := getLogWriter(initLog.LogSavePath + initLog.HighLevelFile + initLog.LogFileExt)

	// 获取日志编码器。该编码器用于格式化日志消息。
	encoder := getEncoder()
	// 每个核心由日志编码器、写入器和日志级别启用器组成。
	lowFileCore := zapcore.NewCore(encoder, lowFileWriteSyncer, lowPriority)
	highFileCore := zapcore.NewCore(encoder, highFileWriteSyncer, highPriority)
	cores = append(cores, lowFileCore, highFileCore)

	// 创建一个用于终端输出的核心，并将其添加到 cores 数组中。
	// 如果是debug级别好需要输出到终端
	if level == "debug" {
		// 创建终端输出编码器
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}

	// 将多个核心组合成一个。
	core := zapcore.NewTee(cores...)
	// 使用 core 创建一个新的 Log 实例，并将其作为指针返回。
	// zap.New 函数创建一个新的 zap.Logger，它是一个使用给定核心和选项的日志记录器。
	// zap.AddCaller() 选项添加了调用者的信息到日志中。
	LG := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(LG) // 可以自行添加
	return &Log{LG}        // 增加函数调用信息
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,   // 结尾字符
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 将Level序列化为全大写字符串。例如, InfoLevel被序列化为INFO
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // 格式化时间戳
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 普通的Log Encoder
}

// 配置日志文件的切割和管理行为,创建一个可以将日志写入到指定文件并进行切割管理的 zapcore.WriteSyncer 对象，用于将日志输出到文件中。
func getLogWriter(filename string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{ // 日志切割
		Filename:   filename,
		MaxSize:    initLog.MaxSize,    // M
		MaxBackups: initLog.MaxBackups, // 备份数量
		MaxAge:     initLog.MaxAge,     // 最大备份天数
		Compress:   initLog.Compress,   // 压缩过期日志
	}
	return zapcore.AddSync(lumberJackLogger) //通过 zapcore.AddSync 函数将 lumberJackLogger 转换为 zapcore.WriteSyncer 类型
}
