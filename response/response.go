package response

import "github.com/gin-gonic/gin"

func ErrorResponse(ctx *gin.Context, code int, mess string) {
	ctx.JSON(200, gin.H{
		"code": code,
		"msg":  msg[code],
	})
}
