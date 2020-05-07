package routes

import (
	middleware "coupanda/middlewares"

	"github.com/gin-gonic/gin"
)

// AdvertismentRouter  #router for advertisment
func AdvertismentRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	v1 := r.Group("/v1")
	{

		v1.POST("/register", Register)
		v1.POST("/login", Login)
		v1.GET("/getAdvertisements", GetAdvertisements)
		v1.Use(middleware.TokenAuthMiddleware())
		{
			v1.POST("/createAdvertisemnt", CreateAdvertisement)

		}

	}
	return r
}
