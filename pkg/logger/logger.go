package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(dev bool) error {
	if err := createDirIfNotExists(); err != nil {
		return err
	}

	level := zapcore.InfoLevel
	if dev {
		level = zapcore.DebugLevel
	}

	fileEncoder := getFileEncoder()
	consoleEncoder := getConsoleEncoder()

	file, err := getFileHandle()
	if err != nil {
		return err
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level),
	)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return nil
}

func getFileHandle() (*os.File, error) {
	path, err := getLogDirPath()
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s/Log File %v.txt", path, time.Now().Format("01-02-2006 15:04:05"))
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func getConsoleEncoder() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeTime = zapcore.TimeEncoder(func(t time.Time, e zapcore.PrimitiveArrayEncoder) {
		e.AppendString(t.UTC().Format(time.RFC3339))
	})
	return zapcore.NewConsoleEncoder(cfg)
}

func getFileEncoder() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	cfg.EncodeTime = zapcore.TimeEncoder(func(t time.Time, e zapcore.PrimitiveArrayEncoder) {
		e.AppendString(t.UTC().Format(time.RFC3339Nano))
	})
	return zapcore.NewConsoleEncoder(cfg)
}

func createDirIfNotExists() error {
	dir, err := getLogDirPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, os.ModePerm)
	} else {
		return err
	}
}

func getLogDirPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/logs", path), nil
}
