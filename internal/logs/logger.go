package logs

import (
	"os"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	LogError(prefix string, err error)
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

	var config zap.Config
	if mode == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	ensureLogDirectory("logs")

	config.OutputPaths = []string{
		"stdout",
		"logs/app.log",
	}
	config.ErrorOutputPaths = []string{
		"stderr",
		"logs/app.log",
	}

	logger, err := config.Build()
	if err != nil {
		panic("❌ Failed to initialize zap logger: " + err.Error())
	}

	return &ZapLogger{sugar: logger.Sugar()}
}

func (l *ZapLogger) Info(msg string, args ...interface{})  { l.sugar.Infof(msg, args...) }
func (l *ZapLogger) Error(msg string, args ...interface{}) { l.sugar.Errorf(msg, args...) }
func (l *ZapLogger) Debug(msg string, args ...interface{}) { l.sugar.Debugf(msg, args...) }
func (l *ZapLogger) Fatal(msg string, args ...interface{}) { l.sugar.Fatalf(msg, args...) }

func (l *ZapLogger) LogError(prefix string, err error) {
	if err != nil {
		l.Error("%s: %v", prefix, err)
	}
}

func ensureLogDirectory(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}
}
