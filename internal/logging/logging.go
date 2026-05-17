package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConfigureLogger initializes production logger.
func ConfigureLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))

	l := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(l)

	return l
}

// ConfigureTestLogger initializes testing logger.
func ConfigureTestLogger() *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))

	l := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(l)

	return l
}
