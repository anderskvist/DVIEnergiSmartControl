package dvi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	ini "gopkg.in/ini.v1"
)

// DVILoginGet is a type for defining the login at DVI Energi webservice
type DVILoginGet struct {
	Usermail     string `json:"usermail"`
	Userpassword string `json:"userpassword"`
	Fabnr        int    `json:"fabnr"`
	Get          DVIGet `json:"get"`
}

// DVILoginSet is a type for defining the login at DVI Energi webservice for setting data
type DVILoginSet struct {
	Usermail     string         `json:"usermail"`
	Userpassword string         `json:"userpassword"`
	Fabnr        int            `json:"fabnr"`
	Set          map[string]int `json:"set"`
}

// DVIGet is a type for defining what information to request from DVI Energi webservice
type DVIGet struct {
	Sensor int `json:"sensor"`
	Relay  int `json:"relay"`
	Timer  int `json:"timer"`
}

// DVISet is a type for defining what information to set to DVI Energi webservice
type DVISet struct {
	CH      int `json:"CH"`
	CHCurve int `json:"CHCurve"`
	CHTemp  int `json:"CHTemp"`
}

// DVIResponse contains
type DVIResponse struct {
	Access string            `json:"Access"`
	Fabnr  int               `json:"fabnr"`
	Output DVIResponseOutput `json:"output"`
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

// GetDviData is to get data from DVI
func GetDviData(cfg *ini.File) DVIResponse {

	debug, _ := cfg.Section("main").Key("debug").Bool()
	debug = false
	data := DVILoginGet{
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

	var dviData DVIResponse
	response, err := http.Post("https://ws.dvienergi.com/API/", "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if debug {
			fmt.Println(jsonPrettyPrint(string(maskPassword(string(data)))))
		}

		err := json.Unmarshal(data, &dviData)
		if err != nil {
			panic(err)
		}
	}
	return dviData
}

// SetDVIData is to set data to DVI
func SetDVIData(cfg *ini.File, set map[string]int) {
	debug, _ := cfg.Section("main").Key("debug").Bool()

	data := DVILoginSet{
		Usermail:     cfg.Section("login").Key("usermail").String(),
		Userpassword: cfg.Section("login").Key("userpassword").String(),
		Fabnr:        cfg.Section("login").Key("fabnr").MustInt(),
		Set:          set}

	jsondata, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Could not convert data to json: %s\n", err)
	} else {
		if debug {
			fmt.Println(jsonPrettyPrint(string(maskPassword(string(jsondata)))))
		}
	}

	var dviData DVIResponse
	response, err := http.Post("https://ws.dvienergi.com/API/", "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if debug {
			fmt.Println(jsonPrettyPrint(string(maskPassword(string(data)))))
		}

		err := json.Unmarshal(data, &dviData)
		if err != nil {
			panic(err)
		}
	}

}
