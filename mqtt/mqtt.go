package mqtt

import (
	"encoding/json"
	"os"

	dvi "github.com/anderskvist/DVIEnergiSmartControl/dvi"
	log "github.com/anderskvist/GoHelpers/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	ini "gopkg.in/ini.v1"

	"fmt"
	"net/url"
	"strconv"
	"time"
)

// HAHVAC is a struct to help auto discovery in Home Assistant
type HAHVAC struct {
	Platform                string   `json:"platform"`
	Name                    string   `json:"name"`
	Qos                     int      `json:"qos,omitempty"`
	PayloadOn               int      `json:"payload_on,omitempty"`
	PayloadOff              int      `json:"payload_off,omitempty"`
	PowerCommandTopic       string   `json:"power_command_topic,omitempty"`
	Modes                   []string `json:"modes,omitempty"`
	ModeStateTopic          string   `json:"mode_state_topic,omitempty"`
	ModeStateTemplate       string   `json:"mode_state_template,omitempty"`
	ModeCommandTopic        string   `json:"mode_command_topic,omitempty"`
	CurrentTemperatureTopic string   `json:"current_temperature_topic,omitempty"`
	MinTemp                 int      `json:"min_temp,omitempty"`
	MaxTemp                 int      `json:"max_temp,omitempty"`
	TempStep                int      `json:"temp_step,omitempty"`
	TemperatureCommandTopic string   `json:"temperature_command_topic,omitempty"`
	TemperatureStateTopic   string   `json:"temperature_state_topic,omitifempty"`
	Retain                  bool     `json:"retain,omitempty"`
}

var pubConnection mqtt.Client
var subConnection mqtt.Client

var (
	hotWaterHVAC = HAHVAC{
		Platform:                "mqtt",
		Name:                    "DVI heat pump - hot water 2",
		Qos:                     2,
		TemperatureStateTopic:   "heatpump/Output/Set/VVTemp",
		TemperatureCommandTopic: "heatpump/Input/Set/VVTemp",
		CurrentTemperatureTopic: "heatpump/Output/Sensor/StoragetankHotwater",
		MinTemp:                 45,
		MaxTemp:                 55,
		TempStep:                1,
		Retain:                  true,
		PowerCommandTopic:       "heatpump/Input/Set/VV",
		PayloadOn:               1,
		PayloadOff:              0,
		Modes:                   []string{"Clock", "Constant On", "Constant Off"},
		ModeStateTopic:          "heatpump/Output/Set/VVClock",
		ModeStateTemplate:       "{% set modes = { '0':'Clock', '1':'Constant On',  '2':'Constant Off'} %}{{ modes[value] if value in modes.keys() else 'off' }}",
		ModeCommandTopic:        "heatpump/Input/Set/VVClock",
	}
	centralHeatingHVAC = HAHVAC{
		Platform:                "mqtt",
		Name:                    "DVI heat pump - heating curve 2",
		Qos:                     2,
		TemperatureStateTopic:   "heatpump/Output/Set/CHCurve",
		TemperatureCommandTopic: "heatpump/Input/Set/CHCurve",
		CurrentTemperatureTopic: "heatpump/Output/Sensor/CentralheatingForward",
		MinTemp:                 0,
		MaxTemp:                 20,
		TempStep:                1,
		Retain:                  true,
		Modes:                   []string{"Off", "On"},
		ModeStateTopic:          "heatpump/Output/Set/CH",
		ModeStateTemplate:       "{% set modes = { '0':'Off', '1':'On'} %}{{ modes[value] if value in modes.keys() else 'Off' }}",
		ModeCommandTopic:        "convert/CH",
	}

	MQTTClientID = "DVIEnergiSmartControl" + string(os.Getpid())
)

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	opts.SetCleanSession(true)
	return opts
}

// MonitorMQTT will monitor MQTT for changes
func MonitorMQTT(cfg *ini.File) {
	mqttURL := cfg.Section("mqtt").Key("url").String()
	uri, err := url.Parse(mqttURL)
	if err != nil {
		log.Fatal(err)
	}

	if subConnection == nil {
		subConnection = connect(MQTTClientID+"sub", uri)
		log.Debug("Connecting to MQTT (sub)")
	}

	subConnection.Subscribe("heatpump/Input/#", 0, func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		set := make(map[string]int)

		log.Noticef("[%s] %s\n", topic, string(payload))
		switch topic {
		case "heatpump/Input/Set/CH":
			temp, _ := strconv.Atoi(string(payload))
			set["CH"] = temp
		case "heatpump/Input/Set/CHCurve":
			temp, _ := strconv.ParseFloat(string(payload), 64)
			set["CHCurve"] = int(temp)
		case "heatpump/Input/Set/CHTemp":
			temp, _ := strconv.ParseFloat(string(payload), 64)
			set["CHTemp"] = int(temp)
		case "heatpump/Input/Set/VV":
			temp, _ := strconv.Atoi(string(payload))
			set["VV"] = temp
		case "heatpump/Input/Set/VVClock":
			//set["VVClock"] = HotWaterClockConvertS2I(string(payload))
			temp, _ := strconv.Atoi(string(payload))
			set["VVClock"] = temp
		case "heatpump/Input/Set/VVTemp":
			temp, _ := strconv.ParseFloat(string(payload), 64)
			set["VVTemp"] = int(temp)
		}
		if len(set) > 0 {
			log.Infof("Setting to DVI: %#v", set)
			dvi.SetDVIData(cfg, set)
		}
	})
}

// SendToMQTT will send DVI data to MQTT
func SendToMQTT(cfg *ini.File, dviData dvi.Response) {
	mqttURL := cfg.Section("mqtt").Key("url").String()
	uri, err := url.Parse(mqttURL)
	if err != nil {
		log.Fatal(err)
	}

	if pubConnection == nil {
		pubConnection = connect(MQTTClientID+"pub", uri)
		log.Debug("Connecting to MQTT (pub)")
	}
	pubConnection.Publish("heatpump/Output/Sensor/BrineForward", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.BrineForward))
	pubConnection.Publish("heatpump/Output/Sensor/BrineReturn", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.BrineReturn))
	pubConnection.Publish("heatpump/Output/Sensor/CentralheatingForward", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.CentralheatingForward))
	pubConnection.Publish("heatpump/Output/Sensor/CentralheatingReturn", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.CentralheatingReturn))
	pubConnection.Publish("heatpump/Output/Sensor/StoragetankHotwater", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.StoragetankHotwater))
	pubConnection.Publish("heatpump/Output/Sensor/Roomtemperature", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.Roomtemperature))
	pubConnection.Publish("heatpump/Output/Sensor/StoragetankCentralheating", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.StoragetankCentralheating))
	pubConnection.Publish("heatpump/Output/Sensor/Outsidetemperature", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.Outsidetemperature))
	pubConnection.Publish("heatpump/Output/Sensor/Solarheating", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.Solarheating))
	pubConnection.Publish("heatpump/Output/Sensor/Highpressure", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.Highpressure))
	pubConnection.Publish("heatpump/Output/Sensor/Lowpressure", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.Lowpressure))

	pubConnection.Publish("heatpump/Output/Relay/Relay1", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay1))
	pubConnection.Publish("heatpump/Output/Relay/Relay2", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay2))
	pubConnection.Publish("heatpump/Output/Relay/Relay3", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay3))
	pubConnection.Publish("heatpump/Output/Relay/Relay4", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay4))
	pubConnection.Publish("heatpump/Output/Relay/Relay5", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay5))
	pubConnection.Publish("heatpump/Output/Relay/Relay6", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay6))
	pubConnection.Publish("heatpump/Output/Relay/Relay7", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay7))
	pubConnection.Publish("heatpump/Output/Relay/Relay8", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay8))
	pubConnection.Publish("heatpump/Output/Relay/Relay9", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay9))
	pubConnection.Publish("heatpump/Output/Relay/Relay10", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay10))
	pubConnection.Publish("heatpump/Output/Relay/Relay11", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay11))
	pubConnection.Publish("heatpump/Output/Relay/Relay12", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay12))
	pubConnection.Publish("heatpump/Output/Relay/Relay13", 0, false, fmt.Sprintf("%d", dviData.Output.Relay.Relay13))

	pubConnection.Publish("heatpump/Output/Set/VV", 0, false, fmt.Sprintf("%d", dviData.Output.UserSettings.HotwaterState))
	//pubConnection.Publish("heatpump/Output/Set/VVClock2", 0, false, fmt.Sprintf("%s", HotWaterClockConvertI2S(dviData.Output.UserSettings.HotwaterClock)))
	pubConnection.Publish("heatpump/Output/Set/VVClock", 0, false, fmt.Sprintf("%d", dviData.Output.UserSettings.HotwaterClock))

	pubConnection.Publish("heatpump/Output/Set/VVTemp", 0, false, fmt.Sprintf("%d", dviData.Output.UserSettings.HotwaterTemp))

	pubConnection.Publish("heatpump/Output/Set/CH", 0, false, fmt.Sprintf("%d", dviData.Output.UserSettings.CentralheatState))
	pubConnection.Publish("heatpump/Output/Set/CHCurve", 0, false, fmt.Sprintf("%d", dviData.Output.UserSettings.CentralheatCurve))
	pubConnection.Publish("heatpump/Output/Set/CHTemp", 0, false, fmt.Sprintf("%d", dviData.Output.UserSettings.CentralheatTemp))
}

func HotWaterClockConvertI2S(i int) string {
	switch i {
	case 0:
		return "Clock"
	case 1:
		return "Constant On"
	case 2:
		return "Constant Off"
	}
	return "Constant On"
}

func HotWaterClockConvertS2I(s string) int {
	switch s {
	case "Clock":
		return 0
	case "Constant On":
		return 1
	case "Constant Off":
		return 2
	}
	return 2
}

// HomeAssistantAutoDiscovery will send DVI data to MQTT
func HomeAssistantAutoDiscovery(cfg *ini.File) {
	mqttURL := cfg.Section("mqtt").Key("url").String()
	uri, err := url.Parse(mqttURL)
	if err != nil {
		log.Fatal(err)
	}

	if pubConnection == nil {
		pubConnection = connect(MQTTClientID+"pub", uri)
		log.Debug("Connecting to MQTT (pub)")
	}

	hwHVAC, _ := json.Marshal(hotWaterHVAC)
	log.Debug(string(hwHVAC))
	pubConnection.Publish("homeassistant/climate/DVIEnergiSmartControl/HotWaterHVAC/config", 1, true, hwHVAC)

	chHVAC, _ := json.Marshal(centralHeatingHVAC)
	log.Debug(string(chHVAC))
	pubConnection.Publish("homeassistant/climate/DVIEnergiSmartControl/CentralHeatingHVAC/config", 1, true, chHVAC)
}
