package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type SensorContent struct {
	Title     string      `json:"title,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

var (
	fileFlag string
	column   []string
)

func init() {
	now := time.Now().Format("01-02 15:04:05")
	flag.StringVar(&fileFlag, "f", now+".csv", "set the output file")
}

func main() {
	flag.Parse()
	r := gin.Default()
	r.POST("/api/sensors", sensorHandler)
	log.Fatalln(r.Run(":8080"))
}

func sensorHandler(c *gin.Context) {
	var sensors map[string]SensorContent
	if err := c.ShouldBindJSON(&sensors); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InternalServerError"})
		return
	}

	file, err := os.OpenFile(fileFlag, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InternalServerError"})
		return
	}
	defer file.Close()

	columnLength := len(column)
	if columnLength == 0 {
		CreateColum(file, sensors)
	}

	var records []string
	for _, c := range column {
		records = append(records, fmt.Sprint(sensors[c].Value))
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(records)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"len":     columnLength,
	})
}

func CreateColum(file io.Writer, sensors map[string]SensorContent) error {
	for key := range sensors {
		column = append(column, key)
	}
	sort.Strings(column)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	return writer.Write(column)
}
