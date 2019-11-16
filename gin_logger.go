package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/yanjiaqi03/golang-logger/utils"
)

const (
	KEY_ENTRY = "key_log_entry"
	KEY_TRACE_ID = "gin_trace_id"
)

func Context(c *gin.Context) ( *logrus.Entry) {
	entry, _ := c.Get(KEY_ENTRY)
	if entry == nil {
		entry = Instance().WithField(KEY_TRACE_ID, utils.UniqueId())
		c.Set(KEY_ENTRY, entry)
	}
	return entry.(*logrus.Entry)
}
