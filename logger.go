package logger

import (
	"github.com/sirupsen/logrus"
	"sync"
	"github.com/lestrrat-go/file-rotatelogs"
	"time"
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/olivere/elastic"
	"github.com/yanjiaqi03/golang-logger/hook"
	"github.com/yanjiaqi03/golang-logger/utils"
)

var logger *logrus.Logger
var once sync.Once

type Builder struct {
	Name string
	LogLevel logrus.Level
	LogPath  string
	ElasticHost string
	LocalHost string
}

func NewBuilder() (*Builder) {
	return &Builder{
		Name: "mylog",
		LogLevel: logrus.DebugLevel,
		LogPath:  "",
		ElasticHost: "http://localhost:9200",
		LocalHost: utils.GetLocalAddress(),
	}
}

func (builder *Builder) SetName(name string) *Builder {
	if len(name) > 0 {
		builder.Name = name
	}
	return builder
}

func (builder *Builder) SetLocalHost(localHost string) *Builder {
	builder.LocalHost = localHost
	return builder
}

func (builder *Builder) SetLevel(level logrus.Level) *Builder {
	builder.LogLevel = level
	return builder
}

func (builder *Builder) SetLogPath(path string) *Builder {
	builder.LogPath = path
	return builder
}

func (builder *Builder) SetElasticHost(host string) *Builder {
	builder.ElasticHost = host
	return builder
}

func (builder *Builder) Build() (*logrus.Logger) {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: utils.TimeFormat,
	})
	logger.SetLevel(builder.LogLevel)
	if len(builder.LogPath) > 0 {
		path := builder.LogPath + builder.Name + "/%Y%m%d/log"
		writer, err := rotatelogs.New(
			path+".%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(-1),
			rotatelogs.WithRotationTime(time.Hour),
		)
		if err != nil {
			logger.Panic(fmt.Sprint("log initial error: ", err))
		}
		logger.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.PanicLevel: writer,
				logrus.FatalLevel: writer,
			},
			&logrus.TextFormatter{
				TimestampFormat: utils.TimeFormat,
			},
		))
	}
	if len(builder.ElasticHost) > 0 {
		// elasticHook
		client, err := elastic.NewClient(elastic.SetURL(builder.ElasticHost))
		if err != nil {
			logger.Panic(err)
		}
		elasticHook := hook.NewElasticHook(client, builder.LocalHost, builder.LogLevel, builder.Name)
		if err != nil {
			logger.Panic(err)
		}
		logger.AddHook(elasticHook)
	}
	return logger
}

func Instance() (*logrus.Logger) {
	if logger == nil {
		once.Do(func() {
			logger = NewBuilder().Build()
		})
	}
	return logger
}