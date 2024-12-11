package user

import "github.com/gin-gonic/gin"

type ProductRouter struct{}

func (pr *ProductRouter) InitProductRouter(Router *gin.RouterGroup) {
	//public router
	productRouterPublic := Router.Group("/product")
	{
		productRouterPublic.GET("/search")
		productRouterPublic.GET("/detail/:id")
	}

	//private router
	productRouterPrivate := Router.Group("/product")
	//productRouterPrivate.Use(middlewares.AuthenMiddleware())
	//productRouterPrivate.Use(middlewares.AuthorizeMiddleware())
	{
		productRouterPrivate.GET("/cart/list")
	}
}
