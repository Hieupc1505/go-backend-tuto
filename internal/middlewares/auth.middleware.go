package middlewares

import (
	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/response"
)

func AuthenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token != "valid-token" {
			response.ErrorResponse(ctx, response.ErrInvalidToke, "")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
