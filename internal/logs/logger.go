package logs

import (
	"go.uber.org/zap"
	"os"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})
}

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

var (
	instance *ZapLogger
	once     sync.Once
)

func Init() {
	once.Do(func() {
		instance = newZapLogger()
	})
}

func Get() *ZapLogger {
	Init()
	return instance
}

func newZapLogger() *ZapLogger {
	mode := viper.GetString("MODE")
	if mode == "" {
		mode = "development"
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:    "timestamp",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "caller",
		MessageKey: "message",
		//StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	level := zap.DebugLevel
	if mode == "production" {
		level = zap.InfoLevel
	}

	stdoutEncoder := WrapEncoderAsPretty(zapcore.NewJSONEncoder(encoderCfg))
	stdoutCore := zapcore.NewCore(stdoutEncoder, zapcore.AddSync(os.Stdout), level)

	fileEncoderCfg := encoderCfg
	fileEncoderCfg.StacktraceKey = "stacktrace"

	var fileEncoder zapcore.Encoder
	var loggingLevel zapcore.Level
	if mode == "production" {
		fileEncoder = zapcore.NewJSONEncoder(fileEncoderCfg)
		loggingLevel = zap.ErrorLevel
	} else {
		fileEncoder = WrapEncoderAsPretty(zapcore.NewJSONEncoder(fileEncoderCfg))
		loggingLevel = zap.DebugLevel
	}

	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(getLogFileWriter()), level)

	core := zapcore.NewTee(stdoutCore, fileCore)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(loggingLevel))
	return &ZapLogger{sugar: logger.Sugar()}
}

func getLogFileWriter() zapcore.WriteSyncer {
	file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}
	return zapcore.AddSync(file)
}

//// Public API
//func Info(msg string, args ...interface{}) {
//	Get().sugar.Info(append([]interface{}{msg}, args...)...)
//}
//
//func Error(args ...interface{}) {
//	//Get().sugar.Error(append([]interface{}{msg}, args...)...)
//	Get().sugar.Errorw("test1", args...)
//
//}
//func Debug(args ...interface{}) {
//	Get().sugar.Debug(args...)
//}
//func Fatal(msg string, args ...interface{}) {
//	Get().sugar.Fatal(args...)
//}
//func Warn(msg string, args ...interface{}) {
//	Get().sugar.Warn(args...)
//}

func (l *ZapLogger) Info(args ...interface{})  { l.sugar.Info(args...) }
func (l *ZapLogger) Error(args ...interface{}) { l.sugar.Errorw("", args) }
func (l *ZapLogger) Debug(args ...interface{}) { l.sugar.Debug(args...) }
func (l *ZapLogger) Fatal(args ...interface{}) { l.sugar.Fatal(args...) }
