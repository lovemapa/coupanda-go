package main

import (
	routes "coupanda/advertiseroutes"

	"github.com/gin-gonic/gin"
)

func router() {

	r := gin.Default()

	r = routes.AdvertismentRouter()

	r.Run(":8080")

}
