package influx

import (
	"strconv"
	"time"

	log "github.com/anderskvist/DVIEnergiSmartControl/log"

	client "github.com/influxdata/influxdb/client/v2"

	"github.com/anderskvist/DVIEnergiSmartControl/dvi"
	ini "gopkg.in/ini.v1"
)

func SaveToInflux(cfg *ini.File, dviData dvi.DVIResponse) {

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
		log.Fatal(err)
	}

	timer_points, err := client.NewPoint(
		"timer",
		tags,
		timers,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}

	relay_points, err := client.NewPoint(
		"relay",
		tags,
		relays,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
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
