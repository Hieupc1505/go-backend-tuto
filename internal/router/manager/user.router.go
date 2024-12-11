package manager

import "github.com/gin-gonic/gin"

type UserManagerRouter struct{}

func (pr *UserManagerRouter) InitUserManagerRouter(Router *gin.RouterGroup) {

	adminRouterPrivate := Router.Group("/admin/user")
	{
		adminRouterPrivate.POST("/active_user")
	}
}
