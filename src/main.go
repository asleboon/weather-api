package main

import (
	"weather-api/src/handlers"

	"github.com/gin-gonic/gin"
)

// TODO: get location by zipCode

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// Set trusted proxies
	r.SetTrustedProxies(nil)
	// load html template files
	r.LoadHTMLGlob("src/templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
		return
	})
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/weather", handlers.WeatherHandler)
		apiV1.GET("/swagger", func(c *gin.Context) {
			c.HTML(200, "swagger.html", nil)
			return
		})
	}

	r.Run(":8080")
}
