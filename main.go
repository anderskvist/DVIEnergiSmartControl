package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	ini "gopkg.in/ini.v1"
)

var debug = true

// DVILogin is a type for defining the login at DVI Energi webservice
type DVILogin struct {
	Usermail     string `json:"usermail"`
	Userpassword string `json:"userpassword"`
	Fabnr        int    `json:"fabnr"`
	Get          DVIGet `json:"get"`
}

// DVIGet is a type for defining what information to request from DVI Energi webservice
type DVIGet struct {
	Sensor int `json:"sensor"`
	Relay  int `json:"relay"`
	Timer  int `json:"timer"`
}

// DVIResponse contains
type DVIResponse struct {
	Access string            `json:"Access"`
	Fabnr  int               `json:"fabnr"`
	Output DVIResponseOutput `json:"output"`
	//	Output []struct {
	//	name string
	//	}
}

// DVIResponseOutput contains
type DVIResponseOutput struct {
	Sensor DVIResponseOutputSensor `json:"sensor"`
	Relay  DVIResponseOutputRelay  `json:"relay"`
	Timer  DVIResponseOutputTimer  `json:"timer"`
}

// DVIResponseOutputSensor contains sensor data
type DVIResponseOutputSensor struct {
	SensorDate                string  `json:"Sensor.Date"`
	CentralheatingForward     float32 `json:"Centralheating.Forward,string"`
	CentralheatingReturn      float32 `json:"Centralheating.Return,string"`
	StoragetankHotwater       float32 `json:"Storagetank.Hotwater,string"`
	Roomtemperature           float32 `json:"Roomtemperature,string"`
	StoragetankCentralheating float32 `json:"Storagetank.Centralheating,string"`
	LVEvaporator1             float32 `json:"LV.Evaporator1,string"`
	Outsidetemperature        float32 `json:"Outsidetemperature,string"`
	Energycatcher             float32 `json:"Energycatcher,string"`
	Solarheating              float32 `json:"Solarheating,string"`
	LVEvaporator2             float32 `json:"LV.Evaporator2,string"`
	Highpressure              float32 `json:"Highpressure,string"`
	Lowpressure               float32 `json:"Lowpressure,string"`
	BrineReturn               float32 `json:"Brine.Return,string"`
	BrineForward              float32 `json:"Brine.Forward,string"`
	HeatmeterFlow             float32 `json:"Heatmeter.Flow,string"`
	HeatmeterKW               float32 `json:"Heatmeter.kW,string"`
	HeatmeterForward          float32 `json:"Heatmeter.Forward,string"`
	HeatmeterReturn           float32 `json:"Heatmeter.Return,string"`
	HeatmeterKWh              float32 `json:"Heatmeter.kWh,string"`
	PowermeterKW              float32 `json:"Powermeter.kW,string"`
	PowermeterKWh             float32 `json:"Powermeter.kWh,string"`
}

// DVIResponseOutputRelay contains relay data
type DVIResponseOutputRelay struct {
	Relay1  int `json:"Relay1,string"`
	Relay2  int `json:"Relay2,string"`
	Relay3  int `json:"Relay3,string"`
	Relay4  int `json:"Relay4,string"`
	Relay5  int `json:"Relay5,string"`
	Relay6  int `json:"Relay6,string"`
	Relay7  int `json:"Relay7,string"`
	Relay8  int `json:"Relay8,string"`
	Relay9  int `json:"Relay9,string"`
	Relay10 int `json:"Relay10,string"`
	Relay11 int `json:"Relay11,string"`
	Relay12 int `json:"Relay12,string"`
	Relay13 int `json:"Relay13,string"`
	Relay14 int `json:"Relay14,string"`
}

// DVIResponseOutputTimer contains strings, but currently data is integer - this may be changed to floating point numbers later
type DVIResponseOutputTimer struct {
	Compressor    int `json:"compressor,string"`
	Warmwater     int `json:"warmwater,string"`
	Pluswarm      int `json:"pluswarm,string"`
	Energicapture int `json:"energicapture,string"`
	Sunwarm       int `json:"sunwarm,string"`
	Suntoearth    int `json:"suntoearth,string"`
	Cooling       int `json:"cooling,string"`
}

// maskPassword will find password in the json string and mask it
func maskPassword(json string) string {
	var regex = regexp.MustCompile(`userpassword\":([ ]?)\"([a-zA-Z0-9]+)\"`)
	s := regex.ReplaceAllString(json, `userpassword":$1"********"`)
	return string(s)
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	data := DVILogin{
		Usermail:     cfg.Section("login").Key("usermail").String(),
		Userpassword: cfg.Section("login").Key("userpassword").String(),
		Fabnr:        cfg.Section("login").Key("fabnr").MustInt(),
		Get: DVIGet{
			Sensor: cfg.Section("get").Key("sensor").MustInt(0),
			Relay:  cfg.Section("get").Key("relay").MustInt(0),
			Timer:  cfg.Section("get").Key("timer").MustInt(0)}}

	jsondata, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Could not convert data to json: %s\n", err)
	} else {
		if debug {
			fmt.Println(jsonPrettyPrint(string(maskPassword(string(jsondata)))))
		}
	}

	response, err := http.Post("https://ws.dvienergi.com/API/", "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if debug {
			fmt.Println(jsonPrettyPrint(string(maskPassword(string(data)))))
		}

		var dviData DVIResponse
		err := json.Unmarshal(data, &dviData)
		if err != nil {
			panic(err)
		}

		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     cfg.Section("influxdb").Key("url").String(),
			Username: cfg.Section("influxdb").Key("username").String(),
			Password: cfg.Section("influxdb").Key("password").String(),
		})

		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  cfg.Section("influxdb").Key("database").String(),
			Precision: "s",
		})

		tags := map[string]string{"fabnr": strconv.Itoa(dviData.Fabnr)}

		sensors := map[string]interface{}{
			"BrineForward":              dviData.Output.Sensor.BrineForward,
			"BrineReturn":               dviData.Output.Sensor.BrineReturn,
			"CentralheatingForward":     dviData.Output.Sensor.CentralheatingForward,
			"CentralheatingReturn":      dviData.Output.Sensor.CentralheatingReturn,
			"Energycatcher":             dviData.Output.Sensor.Energycatcher,
			"HeatmeterFlow":             dviData.Output.Sensor.HeatmeterFlow,
			"HeatmeterForward":          dviData.Output.Sensor.HeatmeterForward,
			"HeatmeterKW":               dviData.Output.Sensor.HeatmeterKW,
			"HeatmeterKWh":              dviData.Output.Sensor.HeatmeterKWh,
			"HeatmeterReturn":           dviData.Output.Sensor.HeatmeterReturn,
			"Highpressure":              dviData.Output.Sensor.Highpressure,
			"LVEvaporator1":             dviData.Output.Sensor.LVEvaporator1,
			"LVEvaporator2":             dviData.Output.Sensor.LVEvaporator2,
			"Lowpressure":               dviData.Output.Sensor.Lowpressure,
			"Outsidetemperature":        dviData.Output.Sensor.Outsidetemperature,
			"PowermeterKW":              dviData.Output.Sensor.PowermeterKW,
			"PowermeterKWh":             dviData.Output.Sensor.PowermeterKWh,
			"Roomtemperature":           dviData.Output.Sensor.Roomtemperature,
			"Solarheating":              dviData.Output.Sensor.Solarheating,
			"StoragetankCentralheating": dviData.Output.Sensor.StoragetankCentralheating,
			"StoragetankHotwater":       dviData.Output.Sensor.StoragetankHotwater,
		}

		timers := map[string]interface{}{
			"Compressor":    dviData.Output.Timer.Compressor,
			"Cooling":       dviData.Output.Timer.Cooling,
			"Energicapture": dviData.Output.Timer.Energicapture,
			"Pluswarm":      dviData.Output.Timer.Pluswarm,
			"Suntoearth":    dviData.Output.Timer.Suntoearth,
			"Sunwarm":       dviData.Output.Timer.Sunwarm,
			"Warmwater":     dviData.Output.Timer.Warmwater,
		}

		relays := map[string]interface{}{
			"Relay1":  dviData.Output.Relay.Relay1,
			"Relay2":  dviData.Output.Relay.Relay2,
			"Relay3":  dviData.Output.Relay.Relay3,
			"Relay4":  dviData.Output.Relay.Relay4,
			"Relay5":  dviData.Output.Relay.Relay5,
			"Relay6":  dviData.Output.Relay.Relay6,
			"Relay7":  dviData.Output.Relay.Relay7,
			"Relay8":  dviData.Output.Relay.Relay8,
			"Relay9":  dviData.Output.Relay.Relay9,
			"Relay10": dviData.Output.Relay.Relay10,
			"Relay11": dviData.Output.Relay.Relay11,
			"Relay12": dviData.Output.Relay.Relay12,
			"Relay13": dviData.Output.Relay.Relay13,
			"Relay14": dviData.Output.Relay.Relay14,
		}

		sensor_points, err := client.NewPoint(
			"sensor",
			tags,
			sensors,
			time.Now(),
		)
		if err != nil {
			log.Fatalln("Error: ", err)
		}

		timer_points, err := client.NewPoint(
			"timer",
			tags,
			timers,
			time.Now(),
		)
		if err != nil {
			log.Fatalln("Error: ", err)
		}

		relay_points, err := client.NewPoint(
			"relay",
			tags,
			relays,
			time.Now(),
		)
		if err != nil {
			log.Fatalln("Error: ", err)
		}

		bp.AddPoint(sensor_points)
		bp.AddPoint(timer_points)
		bp.AddPoint(relay_points)

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		// Close client resources
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}

	}

}
