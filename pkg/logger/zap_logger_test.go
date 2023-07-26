package logger

import (
	"testing"

	"go.uber.org/zap"
)

var (
	logger *InitStruct
	LG     *Log
)

const level string = "debug"

func TestNewLogger(t *testing.T) {
	logger = &InitStruct{
		LogSavePath:   "./storage/appLogs/",
		LogFileExt:    ".log",
		MaxSize:       10,
		MaxBackups:    7,
		MaxAge:        30,
		Compress:      false,
		LowLevelFile:  "info",
		HighLevelFile: "error",
	}
	LG = NewLogger(logger, level)
	zap.ReplaceGlobals(LG.Logger)
	zap.S().Error("hahahh")
	zap.L().Error("xixixi")
	zap.S().Info("start!!")
}
