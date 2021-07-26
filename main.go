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

	go func() {
		for {
			dataRaw, err := waterlevel.Fetch()
			if err != nil {
				log.Println(err)
			}

			dataParsed, err := waterlevel.Parse(dataRaw)
			if err != nil {
				log.Println(err)
			}

			dataByStation, err = waterlevel.ParseByLocation(dataParsed)
			if err != nil {
				log.Println(err)
			}

			last := dataParsed[len(dataParsed)-1]
			log.Printf("Number of datasets: %d", len(dataParsed))
			log.Printf("Last dataset: %s %s %v", last.StationName, last.Time, last.Value)

			time.Sleep(5 * time.Minute)

		}
	}()

	stations := []string{}
	for k := range dataByStation {
		stations = append(stations, k)
	}

	sort.Strings(stations)

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
