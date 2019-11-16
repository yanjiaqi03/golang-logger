package hook

import (
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"context"
	"time"
	"strings"
	"logger/utils"
)

type ElasticHook struct {
	client *elastic.Client
	host   string
	index  string
	levels []logrus.Level
}

type message struct {
	Host      string
	Timestamp string
	Message   string
	Data      logrus.Fields
	Level     string
}

func NewElasticHook(client *elastic.Client, host string, level logrus.Level, index string) *ElasticHook {
	var levels []logrus.Level
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}

	return &ElasticHook{
		client: client,
		host:   host,
		index:  index,
		levels: levels,
	}
}

func (hook *ElasticHook) Fire(entry *logrus.Entry) error {
	msg := message{
		Host:      hook.host,
		Timestamp: time.Now().Format(utils.TimeFormat),
		Level:     strings.ToUpper(entry.Level.String()),
		Message:   entry.Message,
		Data:      entry.Data,
	}

	_, err := hook.client.
		Index().
		Index(hook.index).
		BodyJson(msg).
		Do(context.Background())

	return err
}

func (hook *ElasticHook) Levels() []logrus.Level {
	return hook.levels
}
