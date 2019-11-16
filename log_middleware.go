package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/url"
	"time"
	"strconv"
	"github.com/sirupsen/logrus"
)

const (
	KEY_PATH         = "gin_path"
	KEY_QUERY        = "gin_query"
	KEY_REQUEST_BODY = "gin_request_body"
	KEY_COST_TIME    = "gin_cost_time"
	KEY_RESPONSE     = "gin_response"
)

func GinLogMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// 公参
		entry := Context(c)
		var bodyStr url.Values
		if c.Request.Body != nil {
			bodyBytes, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				entry.Error(err.Error())
			} else {
				bodyStr, err = url.ParseQuery(string(bodyBytes))
				if err != nil {
					entry.Error(err.Error())
				}
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
			}
		}
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		entry.WithFields(logrus.Fields{
			KEY_PATH:         c.Request.URL.Path,
			KEY_QUERY:        c.Request.URL.Query().Encode(),
			KEY_REQUEST_BODY: bodyStr.Encode(),
		}).Info("")

		defer TryError(c)

		c.Next()

		costTimeStr := strconv.FormatInt(time.Now().UnixNano()/1e6-startTime.UnixNano()/1e6, 10) + "ms"
		entry.WithFields(logrus.Fields{
			KEY_RESPONSE:  blw.body.String(),
			KEY_COST_TIME: costTimeStr,
		}).Info("")
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
