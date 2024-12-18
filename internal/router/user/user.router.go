package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/middlewares"
	"hieupc05.github/backend-server/internal/wire"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// ur := repos.NewUserRepository()
	// us := services.NewUserServices(ur)
	// userHanlerNonDependency := controllers.NewUserController(us)

	//TODO: WIRE go
	//Dependency injection (DI) DI java
	userHanlerNonDependency, err := wire.InitUserRouterHandler(global.Config.Token.SecretKey, global.TokenMaker)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	UserRouterPublic := Router.Group("/user")
	{
		UserRouterPublic.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		UserRouterPublic.POST("/register", userHanlerNonDependency.Register)
		UserRouterPublic.POST("/verify_otp", userHanlerNonDependency.CreateUser)
		UserRouterPublic.POST("/login", userHanlerNonDependency.Login)
	}
	//private routercd
	UserRouterPrivate := Router.Group("/user")
	// //UserRouterPrivate.Use(middlewares.Limit())
	UserRouterPrivate.Use(middlewares.AuthenMiddleware(global.TokenMaker))
	// //UserRouterPrivate.Use(middlewares.Permission())
	{
		UserRouterPrivate.GET("/get_info", userHanlerNonDependency.GetUserInfo)
	}

	// fmt.Printf("PublicRouter: %+v\n", UserRouterPublic)
	// fmt.Printf("PrivateRouter: %+v\n", UserRouterPrivate)
	//public router

}
