package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/internal/utils/token"
	"hieupc05.github/backend-server/response"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthenMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			rsp := response.ErrorResponse(response.ErrUnauthorizedInvalidToken)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rsp)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			rsp := response.ErrorResponse(response.ErrUnauthorizedInvalidToken)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rsp)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			rsp := response.ErrorResponse(response.ErrUnauthorizedInvalidToken)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rsp)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			rsp := response.ErrorResponse(response.ErrUnauthorizedInvalidToken)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rsp)
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
