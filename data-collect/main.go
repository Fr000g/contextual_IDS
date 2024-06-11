package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mu       sync.Mutex
	fileName string
	topic    string
	broker   string
	port     int
	output   string
)

var columns = []string{
	"time", "temperature_dht", "temperature_bmp", "humidity", "light", "sound", "pressure", "movement", "vibration", "accX", "accY", "accZ",
}

type EMS struct {
	Time           string  `json:"time"`
	TemperatureDht int     `json:"temperature_dht"`
	TemperatureBmp float64 `json:"temperature_bmp"`
	Humidity       int     `json:"humidity"`
	Light          int     `json:"light"`
	Sound          int     `json:"sound"`
	Pressure       int     `json:"pressure"`
	Movement       int     `json:"movement"`
	Vibration      int     `json:"vibration"`
	AccX           float64 `json:"accX"`
	AccY           float64 `json:"accY"`
	AccZ           float64 `json:"accZ"`
}

func parseFlags() {
	flag.StringVar(&topic, "t", "/sensors", "topic for mqtt")
	flag.StringVar(&broker, "b", "localhost", "broker for mqtt")
	flag.IntVar(&port, "p", 1883, "port for broker")
	flag.StringVar(&output, "o", "", "output file")
	flag.Parse()
}

func main() {
	parseFlags()
	fmt.Println("Start with topic:", topic)

	if output == "" {
		fileName = GenCSVFileName()
	} else {
		fileName = output
	}
	fmt.Println("Save file to:", fileName)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = onConnectHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Failed to connect:", token.Error())
		os.Exit(1)
	}
	defer client.Disconnect(250)

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println("Failed to subscribe:", token.Error())
		os.Exit(1)
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)

	// kepp running
	select {}

}
func messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("message received: %s from %s\n", msg.Payload(), msg.Topic())
	mu.Lock()
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()
		CreateColumn(file, msg)
	} else {
		file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
		defer file.Close()
		WriteValue(file, msg)
	}
	mu.Unlock()
}

func onConnectHandler(mqtt.Client) {
	fmt.Println("connected!")
}

func GenCSVFileName() string {
	s := time.Now().Format("2006-1-2_15:04:05") + ".csv"
	return s
}

func WriteValue(file *os.File, msg mqtt.Message) {
	var ems EMS
	err := json.Unmarshal(msg.Payload(), &ems)
	if err != nil {
		return
	}

	ems.Time = time.Now().Format("2006-01-02 15:04:05")

	values := []string{
		ems.Time,
		fmt.Sprintf("%d", ems.TemperatureDht),
		fmt.Sprintf("%.2f", ems.TemperatureBmp),
		fmt.Sprintf("%d", ems.Humidity),
		fmt.Sprintf("%d", ems.Light),
		fmt.Sprintf("%d", ems.Sound),
		fmt.Sprintf("%d", ems.Pressure),
		fmt.Sprintf("%d", ems.Movement),
		fmt.Sprintf("%d", ems.Vibration),
		fmt.Sprintf("%.2f", ems.AccX),
		fmt.Sprintf("%.2f", ems.AccY),
		fmt.Sprintf("%.2f", ems.AccZ),
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(values)
}

func CreateColumn(file *os.File, msg mqtt.Message) {
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(columns)
}
