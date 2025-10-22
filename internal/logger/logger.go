package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() *zap.SugaredLogger {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// Поток stdout для Info и ниже
	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l < zapcore.ErrorLevel
	})

	// Поток stderr для Error и выше
	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), infoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), errorLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger.Sugar()
}
