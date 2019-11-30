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
		var trace_id = utils.UniqueId()
		entry = Instance().WithField(KEY_TRACE_ID, trace_id)
		c.Set(KEY_ENTRY, entry)
		c.Header(KEY_TRACE_ID, trace_id)
	}
	return entry.(*logrus.Entry)
}
