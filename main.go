package main

import (
	"os"
	"time"

	"github.com/anderskvist/GoHelpers/version"

	"github.com/anderskvist/DVIEnergiSmartControl/dvi"
	"github.com/anderskvist/DVIEnergiSmartControl/influx"
	"github.com/anderskvist/DVIEnergiSmartControl/log"
	"github.com/anderskvist/DVIEnergiSmartControl/mqtt"

	ini "gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load(os.Args[1])

	if err != nil {
		log.Criticalf("Fail to read file: %v", err)
		os.Exit(1)
	}

	log.Infof("GoRenoWeb version: %s.\n", version.Version)

	influxconfig := false
	mqttconfig := false

	if cfg.Section("influxdb").Key("url").String() != "" {
		log.Info("Activating InfluxDB plugin")
		influxconfig = true
	}
	if cfg.Section("mqtt").Key("url").String() != "" {
		log.Info("Activating MQTT plugin")
		mqttconfig = true

		go mqtt.MonitorMQTT(cfg)
	}

	poll := cfg.Section("main").Key("poll").MustInt(60)
	log.Infof("Polltime is %d seconds.\n", poll)

	ticker := time.NewTicker(time.Duration(poll) * time.Second)
	for ; true; <-ticker.C {
		log.Notice("Tick")
		log.Info("Getting data from DVI")
		dviData := dvi.GetDviData(cfg)
		log.Info("Done getting data from DVI")

		if influxconfig {
			log.Info("Saving data to InfluxDB")
			influx.SaveToInflux(cfg, dviData)
			log.Info("Done saving to InfluxDB")
		}

		if mqttconfig {
			log.Info("Sending data to MQTT")
			mqtt.SendToMQTT(cfg, dviData)
			log.Info("Done sending to MQTT")
		}
	}
}
