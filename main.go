package main

import (
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marians/lanuv-nrw-water-level-api/pkg/waterlevel"
)

func main() {
	var dataByStation map[string]waterlevel.Measurement
	var stations []string
	var lastModified string

	go func() {
		for {
			dataRaw, newLastModified, fetchErr := waterlevel.Fetch(lastModified)
			if fetchErr != nil {
				if fetchErr.Error() != "not modified" {
					log.Println(fetchErr)
				}
			} else {

				dataParsed, err := waterlevel.Parse(dataRaw)
				if err != nil {
					log.Println(err)
				} else {
					dataByStation, err = waterlevel.ParseByLocation(dataParsed)
					if err != nil {
						log.Println(err)
					} else {

						for k := range dataByStation {
							stations = append(stations, k)
						}
						sort.Strings(stations)

						log.Printf("Fetched %d datasets published %s", len(dataParsed), newLastModified)

						if newLastModified != "" {
							lastModified = newLastModified
						}
					}
				}
			}

			time.Sleep(5 * time.Minute)

		}
	}()

	router := gin.Default()

	router.GET("/:station", func(c *gin.Context) {
		station := c.Param("station")

		_, ok := dataByStation[station]
		if ok {
			c.JSON(200, gin.H{
				"value": dataByStation[station].Value,
				"time":  dataByStation[station].Time,
			})
		} else {
			c.JSON(404, gin.H{
				"message": "Station not found",
			})
		}
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"stations": stations,
		})
	})

	router.Run()
}
