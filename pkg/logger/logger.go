package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int

const (
	INFO LogLevel = iota
	DEBUG
	WARN
	ERROR
)

// Log will log a message to console and to file. Recommended usage:
//
// logger.Log(logger.INFO, "the message", logger.String("name", "john doe"))
func Log(level LogLevel, msg string, fields ...zap.Field) {
	switch level {
	case DEBUG:
		zap.L().Debug(msg, fields...)
	case INFO:
		zap.L().Info(msg, fields...)
	case WARN:
		zap.L().Warn(msg, fields...)
	case ERROR:
		zap.L().Error(msg, fields...)
	}
}

func Flush() {
	// since we log both to console and file we must sync. we ignore error
	// as it will always complain since console is not syncable
	// https://github.com/uber-go/zap/issues/880
	zap.L().Sync()
}

// InitLogger will initialize a logger that logs messages both to console and to file
// dev will determine the log level. If false, the log level is INFO, otherwise DEBUG.
// if filename is provided it will be used for the name of the log file. otherwise,
// the name will be determined by the current time and date
func InitLogger(dev bool, filename *string) error {
	if err := createDirIfNotExists(); err != nil {
		return err
	}

	level := zapcore.InfoLevel
	if dev {
		level = zapcore.DebugLevel
	}

	fileEncoder := getFileEncoder()
	consoleEncoder := getConsoleEncoder()

	file, err := getFileHandle(filename)
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

func getFileHandle(filenameOverride *string) (*os.File, error) {
	path, err := getLogDirPath()
	if err != nil {
		return nil, err
	}

	var filename string
	if filenameOverride != nil {
		filename = fmt.Sprintf("%s/%s", path, *filenameOverride)
	} else {
		filename = fmt.Sprintf("%s/Log File %v.txt", path, time.Now().Format("01-02-2006 15:04:05"))
	}
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
