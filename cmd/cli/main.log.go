package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// sugar := zap.NewExample().Sugar()
	// sugar.Infof("hello name: %s, age: %d", "hieu", 20)

	// // logger
	// logger := zap.NewExample()
	// logger.Info("Hello Example")

	// //Develop
	// logger, _ = zap.NewDevelopment()
	// logger.Info("Hello newDevelopment")

	// //Production
	// logger, _ = zap.NewProduction()
	// logger.Info("Hello newProduction")

	encoder := getEndcoderLog()
	sync := getWritersync()
	core := zapcore.NewCore(encoder, sync, zapcore.InfoLevel)

	logger := zap.New(core, zap.AddCaller())
	logger.Info("Info log", zap.Int("line", 1))
	logger.Error("Error log", zap.Int("line", 2))

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

func getWritersync() zapcore.WriteSyncer {
	file, err := os.OpenFile("./log/log.txt", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	syncFile := zapcore.AddSync(file)
	syncConsole := zapcore.AddSync(os.Stderr)

	return zapcore.NewMultiWriteSyncer(syncConsole, syncFile)
}
