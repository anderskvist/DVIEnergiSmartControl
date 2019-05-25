package mqtt

import (
	dvi "github.com/anderskvist/DVIEnergiSmartControl/dvi"
	log "github.com/anderskvist/DVIEnergiSmartControl/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	ini "gopkg.in/ini.v1"

	"fmt"
	"net/url"
	"strconv"
	"time"
)

var pubConnection mqtt.Client
var subConnection mqtt.Client

// temporary until we have read access from DVI
var CH int = 1
var CHCurve float64 = 12.0
var CHTemp float64 = 20.0
var VV int = 1
var VVClock int = 0
var VVTemp float64 = 52.0

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

func listen(uri *url.URL, topic string) {
	client := connect("sub", uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

// MonitorMQTT will monitor MQTT for changes
func MonitorMQTT(cfg *ini.File) {
	mqttURL := cfg.Section("mqtt").Key("url").String()
	uri, err := url.Parse(mqttURL)
	if err != nil {
		log.Fatal(err)
	}

	if subConnection == nil {
		subConnection = connect("sub", uri)
	}

	subConnection.Subscribe("heatpump/Input/#", 0, func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		set := make(map[string]int)

		log.Noticef("[%s] %s\n", topic, string(payload))
		switch topic {
		case "heatpump/Input/Set/CH":
			CH, _ = strconv.Atoi(string(payload))
			set["CH"] = CH
		case "heatpump/Input/Set/CHCurve":
			CHCurve, _ = strconv.ParseFloat(string(payload), 64)
			set["CHCurve"] = int(CHCurve)
		case "heatpump/Input/Set/CHTemp":
			CHTemp, _ = strconv.ParseFloat(string(payload), 64)
			set["CHTemp"] = int(CHTemp)
		case "heatpump/Input/Set/VV":
			VV, _ = strconv.Atoi(string(payload))
			set["VV"] = VV
		case "heatpump/Input/Set/VVClock":
			VVClock, _ = strconv.Atoi(string(payload))
			set["VVClock"] = VVClock
		case "heatpump/Input/Set/VVTemp":
			VVTemp, _ = strconv.ParseFloat(string(payload), 64)
			set["VVTemp"] = int(VVTemp)
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
		pubConnection = connect("pub", uri)
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

	pubConnection.Publish("heatpump/Output/Set/VV", 0, false, fmt.Sprintf("%d", VV))
	pubConnection.Publish("heatpump/Output/Set/VVClock", 0, false, fmt.Sprintf("%d", VVClock))
	pubConnection.Publish("heatpump/Output/Set/VVTemp", 0, false, fmt.Sprintf("%f", VVTemp))

	pubConnection.Publish("heatpump/Output/Set/CH", 0, false, fmt.Sprintf("%d", CH))
	pubConnection.Publish("heatpump/Output/Set/CHCurve", 0, false, fmt.Sprintf("%f", CHCurve))
	pubConnection.Publish("heatpump/Output/Set/CHTemp", 0, false, fmt.Sprintf("%f", CHTemp))
}
