package main

import (
	"github.com/gin-gonic/gin"

	routes "coupanda/advertiseroutes"
)

func router() {

	r := gin.Default()

	r.Use(gin.Recovery())
	r = routes.AdvertismentRouter()

	r.Run(":8080")

}
