package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 创建一个 gin.HandlerFunc 类型的中间件,用于记录请求的相关信息
// 该中间件会记录请求的开始时间、请求路径、查询参数、处理请求所耗费的时间,
// 以及一些其他信息如状态码、请求方法、客户端IP、用户代理和错误信息
// 这些信息使用 zap 库进行日志记录,方便后续的日志分析和监控
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始的时间点
		start := time.Now()
		// 获取请求的路径
		path := c.Request.URL.Path
		// 获取请求的查询参数
		query := c.Request.URL.RawQuery

		// 调用 c.Next() 继续处理链中的下一个中间件或处理程序
		c.Next()

		// 计算从请求开始到现在的耗时
		cost := time.Since(start)
		// 使用 zap 库记录日志,包括路径、状态码、请求方式、路径、查询参数、
		// 客户端IP、用户代理、错误信息和请求耗时
		zap.L().Info(path,
			zap.Int("状态码", c.Writer.Status()),
			zap.String("请求方式", c.Request.Method),
			zap.String("路径", path),
			zap.String("参数", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("错误", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("耗时", cost),
		)
	}
}


// Recovery 创建一个 gin.HandlerFunc,用于处理在执行过程中可能发生的 panic
// 参数 stack 控制是否在日志中包含堆栈信息
func Recovery(stack bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 使用 defer 来确保在函数退出时调用 recover,捕获 panic
        defer func() {
            if err := recover(); err != nil {
                // 初始化一个变量来判断是否是由于 TCP 连接问题导致的 panic
                var brokenPipe bool
                // 检查错误是否为 net.OpError 类型,通常是网络操作错误
                if ne, ok := err.(*net.OpError); ok {
                    // 进一步检查错误是否为 os.SyscallError 类型,通常是系统调用错误
                    if se, ok := ne.Err.(*os.SyscallError); ok {
                        // 判断错误信息中是否包含 "broken pipe" 或 "connection reset by peer",表示连接已断开
                        if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
                            brokenPipe = true
                        }
                    }
                }

                // 尝试获取当前请求的详细信息,用于日志记录
                httpRequest, _ := httputil.DumpRequest(c.Request, false)
                // 如果是由于连接断开导致的 panic,则记录 error 级别的日志并终止请求
                if brokenPipe {
                    zap.L().Error(c.Request.URL.Path,
                        zap.Any("错误", err),
                        zap.String("request", string(httpRequest)),
                    )
                    
                    c.Error(err.(error)) 
                    c.Abort()
                    return
                }

                // 根据 stack 参数值,决定是否记录堆栈信息
                if stack {
                    // 如果 stack 为 true,则在日志中包含堆栈信息
                    zap.L().Error("[Recovery from panic]",
                        zap.Any("错误", err),
                        zap.String("request", string(httpRequest)),
                        zap.String("stack", string(debug.Stack())),
                    )
                } else {
                    // 如果 stack 为 false,则不包含堆栈信息
                    zap.L().Error("[Recovery from panic]",
                        zap.Any("错误", err),
                        zap.String("request", string(httpRequest)),
                    )
                }
                // 终止当前请求,并返回 InternalServerError 状态码
                c.AbortWithStatus(http.StatusInternalServerError)
            }
        }()
        // 继续执行后续的 gin.HandlerFunc
        c.Next()
    }
}