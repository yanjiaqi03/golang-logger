package logger

import (
	"github.com/gin-gonic/gin"
	"runtime/debug"
	"net/http"
	"strings"
	"fmt"
)

func TryError(c *gin.Context) {
	r := recover()
	if r == nil {
		return
	}
	entry := Context(c)
	var stack []string
	stack = append(stack, fmt.Sprintf("%v", r))
	stackMsgArr := strings.Split(string(debug.Stack()), "\n")
	for _, item := range stackMsgArr {
		//if strings.Contains(item, "debug/stack.go") || strings.Contains(item, "log_recover.go") || strings.Contains(item, "[running]"){
		//	continue
		//}
		stack = append(stack, strings.Replace(item, "\t", "", -1))
	}
	entry.WithField("stack", stack).Error("")
	//for _, s := range stack {
	//	fmt.Println(s)
	//}
	c.JSON(http.StatusInternalServerError, nil)
}