package kilon

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				log.Printf("%s\n", message) // 打印错误信息
				log.Printf("%s\n", buf[:n]) // 打印堆栈信息
				c.Fail(http.StatusInternalServerError, "Internal Server Error") // 返回服务器错误给客户端
			}
		}()

		c.Next()
	}
}