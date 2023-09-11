package handlers

import (
	"net/http"
	"time"
	"weather-api/src/services"

	"github.com/gin-gonic/gin"
)

var inMemoryCache = services.NewInMemoryCache()

func WeatherHandler(c *gin.Context) {

	// Try to get the data from the cache first
	cacheKey := "some_unique_key_based_on_request_params"
	if data, found := inMemoryCache.Get(cacheKey); found {
		c.JSON(http.StatusOK, gin.H{"data": data})
		return
	}

	// If not in cache, get the location data
	locationData, err := services.GetLocation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Then, use the location data to get the weather data
	weatherData, err := services.GetWeather(locationData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store the weather data in the cache with a 10-minute expiration
	inMemoryCache.Set(cacheKey, weatherData, 10*time.Minute)

	// Finally, respond with the weather data
	c.JSON(http.StatusOK, gin.H{"data": weatherData})
}
