package pkg

import (
	"gin_learn/webook/internal/web/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WrapBody[T any](fn func(req T) (utils.Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(req)
		if err != nil {
			// 处理err
		}
		ctx.JSON(http.StatusOK, res)
	}

}
