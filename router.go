package main

import (
	routes "coupanda/advertiseroutes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func router() {

	r := gin.Default()
	
	r = routes.AdvertismentRouter()
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	r.Run(":" + port)

}
