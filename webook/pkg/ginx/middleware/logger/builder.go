package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

// MiddlewareBuilder 注意点
// 1、小心日志内容过多 URL、请求体、响应体都可能过长
// 2、考虑1的问题以及用户可能换用不同的日志框架。要有足够灵活性
// 3、考虑动态开关、结合监听配置文件、小心并发安全 （原子操作 atomic.NewBool()）
type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	//粗暴的穿一个logger 输出方法不一定是什么
	//logger logger.Logger
	loggerFunc func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}
func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}
func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) >= 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if ctx.Request.Body != nil && b.allowReqBody {
			// body读完就没了
			body, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			// 很消耗CPU和内存
			// 因为会引起复制
			if len(body) >= 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}
		defer func() {
			al.Duration = time.Since(start).String()
			//al.Duration = time.Now().Sub(start)
			b.loggerFunc(ctx, al)
		}()
		// 执行到业务逻辑
		ctx.Next()
		// 不是很优雅的初版实现
		//b.logger.Info()
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(str string) (int, error) {
	w.al.RespBody = str
	return w.ResponseWriter.WriteString(str)
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type AccessLog struct {
	// http请求的方法
	Method string
	// Url
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}
