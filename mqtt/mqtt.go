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

	client := connect("pub", uri)
	client.Publish("heatpump/Output/Sensor/BrineForward", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.BrineForward))
	client.Publish("heatpump/Output/Sensor/BrineReturn", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.BrineReturn))
	client.Publish("heatpump/Output/Sensor/CentralheatingForward", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.CentralheatingForward))
	client.Publish("heatpump/Output/Sensor/CentralheatingReturn", 0, false, fmt.Sprintf("%f", dviData.Output.Sensor.CentralheatingReturn))
	time.Sleep(1 * time.Second)
}
