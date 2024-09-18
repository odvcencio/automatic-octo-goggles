package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/odvcencio/jh_test/clients"
)

func main() {
	r := gin.Default()

	r.GET("/getShortForecast", func(c *gin.Context) {
		lat := c.Query("lat")
		long := c.Query("long")

		latlongRegExp := "^(-?\\d+(?:\\.\\d+)?)"

		validLat, err := regexp.MatchString(latlongRegExp, lat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong",
			})

			log.Println(err)

			return
		}

		validLong, err := regexp.MatchString(latlongRegExp, long)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong",
			})

			log.Println(err)

			return
		}

		if !validLat {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Please enter a valid latitude",
			})

			return
		}

		if !validLong {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Please enter a valid longitude",
			})

			return
		}

		forecast, forecastErr := clients.GetShortForecast(lat, long)
		if forecastErr != nil {
			log.Println(forecastErr.Err)
			returnErrorByCode(c, *forecastErr)

			return
		}

		c.JSON(http.StatusOK, forecast)
	})

	r.Run()
}

func returnErrorByCode(c *gin.Context, err clients.RequestError) {
	switch err.StatusCode {
	case http.StatusBadRequest:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data unavailable for requested point",
		})

		return
	case http.StatusInternalServerError:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})

		return
	}
}
