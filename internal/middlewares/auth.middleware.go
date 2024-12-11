package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/internal/utils/token"
	"hieupc05.github/backend-server/response"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthenMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			rsp := response.ErrorResponse(response.ErrStatusUnauthorized, "", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rsp)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			rsp := response.ErrorResponse(response.ErrStatusUnauthorized, "", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rsp)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			rsp := response.ErrorResponse(response.ErrStatusUnauthorized, "", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rsp)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			rsp := response.ErrorResponse(response.ErrStatusUnauthorized, "", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rsp)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
