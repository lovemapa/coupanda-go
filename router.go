package main

import (
	"github.com/gin-gonic/gin"

	routes "coupanda/advertiseroutes"
)

func router() {

	r := gin.Default()

	r.Use(gin.Recovery())
	r = routes.AdvertismentRouter()
	// api := router.Group("/api/v1")
	// {
	// 	basicAuth := api.Group("/")
	// 	basicAuth.Use(hello())
	// 	{
	// 		basicAuth.POST("/createToken", func(c *gin.Context) {

	// 			t, _ := createToken()

	// 			c.JSON(http.StatusOK, gin.H{
	// 				"token": t,
	// 			})

	// 		})

	// 	}
	// basicAuth.Use(tokenAuthMiddleware())
	// {
	// 	basicAuth.POST("/decodeToken", func(c *gin.Context) {

	// 		c.JSON(http.StatusOK, gin.H{
	// 			"claims": c.Keys["user_id"],
	// 		})

	// 	})
	// }
	// }

	r.Run(":8080")

}
