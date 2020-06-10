package routes

import (
	middleware "coupanda/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdvertismentRouter  #router for advertisment
func AdvertismentRouter() *gin.Engine {
	r := gin.Default()
	r.StaticFS("/static", http.Dir("./static"))
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("templates/*")

	r.GET("/", Welcome)
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
