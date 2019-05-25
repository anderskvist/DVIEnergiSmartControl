package main

import (
	"os"
	"time"

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

	for t := range time.NewTicker(time.Duration(poll) * time.Second).C {
		log.Debug("Tick")
		if t == t {
		}
		log.Debug("Getting data from DVI")
		dviData := dvi.GetDviData(cfg)
		log.Debug("Done getting data from DVI")

		if influxconfig {
			log.Debug("Saving data to InfluxDB")
			influx.SaveToInflux(cfg, dviData)
			log.Debug("Done saving to InfluxDB")
		}

		if mqttconfig {
			log.Debug("Sending data to MQTT")
			mqtt.SendToMQTT(cfg, dviData)
			log.Debug("Done sending to MQTT")
		}
	}
}
