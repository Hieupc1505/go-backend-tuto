package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"hieupc05.github/backend-server/setting"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(config setting.LoggerSetting) *LoggerZap {
	logLevel := config.Log_Level
	//debug -> info -> warn -> error -> fatal -> panic
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoder := getEndcoderLog()
	hook := lumberjack.Logger{
		Filename:   config.File_log_name,
		MaxSize:    config.Max_size, // MB
		MaxBackups: config.Max_backups,
		MaxAge:     config.Max_age,  //days
		Compress:   config.Compress, // disabled by default
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		level)
	// logger := zap.New(core, zap.AddCaller())
	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func getEndcoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()

	//1733062577.3111765 -> 2024-12-01T21:16:17.311+0700
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//ts -> Time
	encodeConfig.TimeKey = "time"

	//from info INFO
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	//"caller": "cli/main.go:18"
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encodeConfig)

}
