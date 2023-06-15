package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

type Data struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

var statusWater string
var statusWind string

func main() {
	tz, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal(err.Error())
	}

	s := gocron.NewScheduler(tz)
	s.Every(15).Seconds().Do(func() {
		data, err := os.Open("data.json")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer data.Close()

		byteValue, _ := io.ReadAll(data)
		var result Data
		err = json.Unmarshal(byteValue, &result)
		if err != nil {
			log.Fatal(err.Error())
		}

		rand.Seed(time.Now().UnixNano())
		result.Water = rand.Intn(100) + 1
		result.Wind = rand.Intn(100) + 1

		file, err := os.Create("data.json")
		if err != nil {
			log.Fatal(err.Error())
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = file.Write(jsonData)
		if err != nil {
			log.Fatal(err.Error())
		}

		if result.Water < 5 {
			statusWater = "aman"
		} else if result.Water >= 6 && result.Water <= 8 {
			statusWater = "siaga"
		} else if result.Water > 8 {
			statusWater = "bahaya"
		}

		if result.Wind < 6 {
			statusWind = "aman"
		} else if result.Wind >= 7 && result.Wind <= 15 {
			statusWind = "siaga"
		} else if result.Wind > 15 {
			statusWind = "bahaya"
		}

		jsonString, _ := json.MarshalIndent(result, "", "    ")

		fmt.Println(string(jsonString))
		fmt.Println("status water : " + statusWater)
		fmt.Println("status wind : " + statusWind)

	})

	s.StartAsync()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/data", func(c *gin.Context) {
		data, err := os.Open("data.json")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer data.Close()

		byteValue, _ := io.ReadAll(data)
		var result Data
		err = json.Unmarshal(byteValue, &result)
		if err != nil {
			log.Fatal(err.Error())
		}

		var statusWater string
		var statusWind string

		if result.Water < 5 {
			statusWater = "aman"
		} else if result.Water >= 6 && result.Water <= 8 {
			statusWater = "siaga"
		} else if result.Water > 8 {
			statusWater = "bahaya"
		}

		if result.Wind < 6 {
			statusWind = "aman"
		} else if result.Wind >= 7 && result.Wind <= 15 {
			statusWind = "siaga"
		} else if result.Wind > 15 {
			statusWind = "bahaya"
		}

		fmt.Println("status water : " + statusWater)
		fmt.Println("status wind : " + statusWind)

		c.JSON(200, result)
	})

	//port := os.Getenv("PORT")
	//router.Run(":" + port)
	router.Run(":8080")
}
