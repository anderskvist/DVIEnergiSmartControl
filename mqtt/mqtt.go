package mqtt

import (
	dvi "github.com/anderskvist/DVIEnergiSmartControl/dvi"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	ini "gopkg.in/ini.v1"

	"fmt"
	"log"
	"net/url"
	"time"
)

var pubConnection mqtt.Client

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
	return opts
}

func listen(uri *url.URL, topic string) {
	client := connect("sub", uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

// SendToMQTT will send DVI data to MQTT
func SendToMQTT(cfg *ini.File, dviData dvi.DVIResponse) {
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

}
