package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type SensorContent struct {
	Id      string      `json:"id"`
	Value   interface{} `json:"v"`
	Quality bool        `json:"q"`
	Time    int64       `json:"t"`
}

type Sensors struct {
	Timestamp int64           `json:"timestamp"`
	Values    []SensorContent `json:"values"`
}

var (
	fileFlag string
	columns  []string
)

func init() {
	now := time.Now().Format("01-02_15-04-05")
	flag.StringVar(&fileFlag, "f", now+".csv", "set the output file")
}

func main() {
	flag.Parse()
	r := gin.Default()
	r.POST("/api/sensors", sensorHandler)
	log.Fatalln(r.Run(":8080"))
}

func sensorHandler(c *gin.Context) {
	var sensors Sensors
	if err := c.ShouldBindJSON(&sensors); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	file, err := os.OpenFile(fileFlag, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	if len(columns) == 0 {
		err := createColumns(file, sensors)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create columns"})
			return
		}
	}

	err = writeRecord(file, sensors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"len":     len(columns),
	})
}

func createColumns(file *os.File, sensors Sensors) error {
	columns = append(columns, "time")
	for _, v := range sensors.Values {
		columns = append(columns, v.Id)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(columns)
}

func writeRecord(file *os.File, sensors Sensors) error {
	record := make([]string, len(columns))

	// Add current time to the first column
	record[0] = time.Now().Format("15:04:05.000")
	for _, value := range sensors.Values {
		for i, col := range columns {
			if col == value.Id {
				record[i] = fmt.Sprintf("%v", value.Value)
				break
			}
		}
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(record)
}
