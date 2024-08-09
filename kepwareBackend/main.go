package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
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
	fileFlag        string
	colMutex        sync.RWMutex
	fileInitialized bool
	fileInitMutex   sync.Mutex
	columns         = []string{
		"time",
		"Factory.High bay warehouse.BB encoder x axis impule 2",
		"Factory.High bay warehouse.HB at rest",
		"Factory.High bay warehouse.HB belt backward",
		"Factory.High bay warehouse.HB belt forward",
		"Factory.High bay warehouse.HB delivered block from pickup",
		"Factory.High bay warehouse.HB done",
		"Factory.High bay warehouse.HB dropoff at belt",
		"Factory.High bay warehouse.HB dropoff block",
		"Factory.High bay warehouse.HB dropoff level down",
		"Factory.High bay warehouse.HB encoder x axis impulse 1",
		"Factory.High bay warehouse.HB encoder y axis impulse 1",
		"Factory.High bay warehouse.HB encoder y axis impulse 2",
		"Factory.High bay warehouse.HB holding block",
		"Factory.High bay warehouse.HB horizontal encoder",
		"Factory.High bay warehouse.HB horizontal encoder value reset",
		"Factory.High bay warehouse.HB in block bay position",
		"Factory.High bay warehouse.HB lever down position",
		"Factory.High bay warehouse.HB lever up position",
		"Factory.High bay warehouse.HB light barrier inside",
		"Factory.High bay warehouse.HB light barrier outside",
		"Factory.High bay warehouse.HB motor cantilever back",
		"Factory.High bay warehouse.HB motor cantilever forward",
		"Factory.High bay warehouse.HB motor x towards belt",
		"Factory.High bay warehouse.HB motor x towards rack",
		"Factory.High bay warehouse.HB motor y axis down",
		"Factory.High bay warehouse.HB motor y axis up",
		"Factory.High bay warehouse.HB move block in",
		"Factory.High bay warehouse.HB move block out",
		"Factory.High bay warehouse.HB move position belt",
		"Factory.High bay warehouse.HB move to belt",
		"Factory.High bay warehouse.HB pockup block",
		"Factory.High bay warehouse.HB processed blocks count",
		"Factory.High bay warehouse.HB put to rest",
		"Factory.High bay warehouse.HB ref switch cantilever backwards",
		"Factory.High bay warehouse.HB ref switch cantilever front",
		"Factory.High bay warehouse.HB ref switch x axis",
		"Factory.High bay warehouse.HB ref switch y axis",
		"Factory.High bay warehouse.HB selected block",
		"Factory.High bay warehouse.HB to blockbay position",
		"Factory.High bay warehouse.HB vertical encoder",
		"Factory.High bay warehouse.HB vertical encoder value reset",
		"Factory.High bay warehouse.HB x axis moving to belt",
		"Factory.High bay warehouse.HB x axis moving to rack",
		"Factory.High bay warehouse.HB x position",
		"Factory.High bay warehouse.HB y position",
		"Factory.High bay warehouse.HB_HOLDING_00_B_BLOCK",
		"Factory.High bay warehouse.HB_HOLDING_01_B_BLOCK",
		"Factory.High bay warehouse.HB_HOLDING_02_B_BLOCK",
		"Factory.High bay warehouse.HB_HOLDING_10_R_BLOCK",
		"Factory.High bay warehouse.HB_HOLDING_11_R_BLOCK",
		"Factory.High bay warehouse.HB_HOLDING_12_R_BLOCK",
		"Factory.High bay warehouse.HB_HOLDING_20_W_BLOCK",
		"Factory.High bay warehouse.HB_HOLDING_21_W_BLOCK",
		"Factory.High bay warehouse.HB_HOLDING_22_W_BLOCK",
		"Factory.High bay warehouse.HB_MTB_TOBELT_POS",
		"Factory.High bay warehouse.HB_MTB_TOREST_POS",
		"Factory.High bay warehouse.HB_WH_POS_00_B",
		"Factory.High bay warehouse.HB_WH_POS_01_B",
		"Factory.High bay warehouse.HB_WH_POS_02_B",
		"Factory.High bay warehouse.HB_WH_POS_10_R",
		"Factory.High bay warehouse.HB_WH_POS_11_R",
		"Factory.High bay warehouse.HB_WH_POS_12_R",
		"Factory.High bay warehouse.HB_WH_POS_20_W",
		"Factory.High bay warehouse.HB_WH_POS_21_W",
		"Factory.High bay warehouse.HB_WH_POS_22_W",
		"Factory.High bay warehouse.Run",
		"Factory.High bay warehouse.Start Button",
		"Factory.High bay warehouse.Stop Button",
		"Factory.High bay warehouse.Trail sensor lower",
		"Factory.High bay warehouse.Trail sensor upper",
		"Factory.Multi Processing Station.MPS Bake Block",
		"Factory.Multi Processing Station.MPS Compressor",
		"Factory.Multi Processing Station.MPS Conveyor forward",
		"Factory.Multi Processing Station.MPS Done",
		"Factory.Multi Processing Station.MPS Done Baking",
		"Factory.Multi Processing Station.MPS Done sawing",
		"Factory.Multi Processing Station.MPS Drop Block",
		"Factory.Multi Processing Station.MPS Light conveyor",
		"Factory.Multi Processing Station.MPS Light oven",
		"Factory.Multi Processing Station.MPS Oven feeder retract",
		"Factory.Multi Processing Station.MPS Oven heating",
		"Factory.Multi Processing Station.MPS Oven lamp",
		"Factory.Multi Processing Station.MPS Over feeder extend",
		"Factory.Multi Processing Station.MPS Process Block",
		"Factory.Multi Processing Station.MPS Ref Turntable belt",
		"Factory.Multi Processing Station.MPS Ref Turntable saw",
		"Factory.Multi Processing Station.MPS Ref Turntable vacuum",
		"Factory.Multi Processing Station.MPS Ref Vacuum oven",
		"Factory.Multi Processing Station.MPS Ref Vacuum turntable",
		"Factory.Multi Processing Station.MPS Ref oven inside",
		"Factory.Multi Processing Station.MPS Ref oven outside",
		"Factory.Multi Processing Station.MPS Run",
		"Factory.Multi Processing Station.MPS Saw motor",
		"Factory.Multi Processing Station.MPS Saw sawing",
		"Factory.Multi Processing Station.MPS Start_Button",
		"Factory.Multi Processing Station.MPS Stop_Button",
		"Factory.Multi Processing Station.MPS Turntable clockwise",
		"Factory.Multi Processing Station.MPS Turntable counterclockwise",
		"Factory.Multi Processing Station.MPS Vacuum Hold Back",
		"Factory.Multi Processing Station.MPS Vacuum Pickup Block",
		"Factory.Multi Processing Station.MPS Vacuum oven",
		"Factory.Multi Processing Station.MPS Vacuum turntable",
		"Factory.Multi Processing Station.MPS Valve Feeder",
		"Factory.Multi Processing Station.MPS Valve Oven door",
		"Factory.Multi Processing Station.MPS Valve lowering",
		"Factory.Multi Processing Station.MPS Valve vacuum",
		"Factory.Multi Processing Station.MPS motor vacuum toward oven to VG",
		"Factory.Multi Processing Station.MPS process block to VG",
		"Factory.Sorting Line.SL Block Detected",
		"Factory.Sorting Line.SL Blue ejector",
		"Factory.Sorting Line.SL Color sensor",
		"Factory.Sorting Line.SL Colour",
		"Factory.Sorting Line.SL Compressor",
		"Factory.Sorting Line.SL Conveyor",
		"Factory.Sorting Line.SL Done",
		"Factory.Sorting Line.SL Eject The Block",
		"Factory.Sorting Line.SL Ejecting",
		"Factory.Sorting Line.SL Light Berrier inlet",
		"Factory.Sorting Line.SL Light behind color sensor",
		"Factory.Sorting Line.SL Light blue",
		"Factory.Sorting Line.SL Light inlet",
		"Factory.Sorting Line.SL Light red",
		"Factory.Sorting Line.SL Light white",
		"Factory.Sorting Line.SL Process Blue Block",
		"Factory.Sorting Line.SL Process Red Block",
		"Factory.Sorting Line.SL Process White Block",
		"Factory.Sorting Line.SL Processing",
		"Factory.Sorting Line.SL RUN state",
		"Factory.Sorting Line.SL Red ejector",
		"Factory.Sorting Line.SL Ref Counter",
		"Factory.Sorting Line.SL Start_button",
		"Factory.Sorting Line.SL Stop_button",
		"Factory.Sorting Line.SL TimerOut",
		"Factory.Sorting Line.SL White ejector",
		"Factory.Sorting Line.SL blue barrier inlet to VG",
		"Factory.Sorting Line.SL done to VG",
		"Factory.Sorting Line.SL red barrier inlet to VG",
		"Factory.Sorting Line.SL white barrier inlet to VG",
		"Factory.Vacuum gripper.VG INIT counter",
		"Factory.Vacuum gripper.VG RUN state",
		"Factory.Vacuum gripper.VG Start Button",
		"Factory.Vacuum gripper.VG Stop Button",
		"Factory.Vacuum gripper.VG at drop off position",
		"Factory.Vacuum gripper.VG at pickup position",
		"Factory.Vacuum gripper.VG compressor",
		"Factory.Vacuum gripper.VG crane at rest",
		"Factory.Vacuum gripper.VG crane at wait",
		"Factory.Vacuum gripper.VG crane in",
		"Factory.Vacuum gripper.VG crane out",
		"Factory.Vacuum gripper.VG crane to rest",
		"Factory.Vacuum gripper.VG crane to wait",
		"Factory.Vacuum gripper.VG cycle x/y",
		"Factory.Vacuum gripper.VG direction clockwise",
		"Factory.Vacuum gripper.VG direction counter clockwise",
		"Factory.Vacuum gripper.VG done moving",
		"Factory.Vacuum gripper.VG dropoff block",
		"Factory.Vacuum gripper.VG dropped block",
		"Factory.Vacuum gripper.VG encoder horizontal axis impulse 1",
		"Factory.Vacuum gripper.VG encoder horizontal axis impulse 2",
		"Factory.Vacuum gripper.VG encoder rotate impulse 1",
		"Factory.Vacuum gripper.VG encoder rotate impulse 2",
		"Factory.Vacuum gripper.VG encoder vertical axis impulse 1",
		"Factory.Vacuum gripper.VG encoder vertical axis impulse 2",
		"Factory.Vacuum gripper.VG holding block",
		"Factory.Vacuum gripper.VG holding blue block",
		"Factory.Vacuum gripper.VG holding red block",
		"Factory.Vacuum gripper.VG holding white block",
		"Factory.Vacuum gripper.VG horizontal encoder",
		"Factory.Vacuum gripper.VG motor horizontal axis backwards",
		"Factory.Vacuum gripper.VG motor horizontal axis forward",
		"Factory.Vacuum gripper.VG motor rotate clockwise",
		"Factory.Vacuum gripper.VG motor rotate counter clockwise",
		"Factory.Vacuum gripper.VG motor vertical axis down",
		"Factory.Vacuum gripper.VG motor vertical axis up",
		"Factory.Vacuum gripper.VG move to hb",
		"Factory.Vacuum gripper.VG move to mps",
		"Factory.Vacuum gripper.VG move to sl blue block",
		"Factory.Vacuum gripper.VG move to sl red block",
		"Factory.Vacuum gripper.VG move to sl white block",
		"Factory.Vacuum gripper.VG moving block",
		"Factory.Vacuum gripper.VG not arriver cntr clkwise",
		"Factory.Vacuum gripper.VG pickup block",
		"Factory.Vacuum gripper.VG position hb",
		"Factory.Vacuum gripper.VG position mps",
		"Factory.Vacuum gripper.VG position sl blue block",
		"Factory.Vacuum gripper.VG position sl red block",
		"Factory.Vacuum gripper.VG position sl white block",
		"Factory.Vacuum gripper.VG ref switch horizontal axis",
		"Factory.Vacuum gripper.VG ref switch rotate",
		"Factory.Vacuum gripper.VG ref switch vertical axis",
		"Factory.Vacuum gripper.VG rotate encoder",
		"Factory.Vacuum gripper.VG rotate motor encoder reset",
		"Factory.Vacuum gripper.VG rotate motor reset",
		"Factory.Vacuum gripper.VG rotation motor encoder value reset",
		"Factory.Vacuum gripper.VG to drop off position",
		"Factory.Vacuum gripper.VG vacuum",
		"Factory.Vacuum gripper.VG vertical encoder",
		"Factory.Vacuum gripper.VG vertical encoder value reset",
		"Factory.Vacuum gripper.VG x out position",
		"Factory.Vacuum gripper.VG y out position",
	}
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

	fileInitMutex.Lock()
	if !fileInitialized {
		// Initialize the file with columns
		if err := initializeFile(); err != nil {
			fileInitMutex.Unlock()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize file"})
			return
		}
		fileInitialized = true
	}
	fileInitMutex.Unlock()

	file, err := os.OpenFile(fileFlag, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

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

func initializeFile() error {
	file, err := os.OpenFile(fileFlag, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if os.IsExist(err) {
		// File already exists, no need to initialize
		return nil
	} else if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(columns)
}

func writeRecord(file *os.File, sensors Sensors) error {
	colMutex.RLock()
	defer colMutex.RUnlock()

	record := make([]string, len(columns))

	// Add current time to the first column
	record[0] = time.Now().Format("2006-01-02 15:04:05.000")
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
