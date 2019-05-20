package main

import (
	"fmt"
	"os"
	"time"

	"github.com/anderskvist/DVIEnergiSmartControl/dvi"
	"github.com/anderskvist/DVIEnergiSmartControl/influx"

	ini "gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load(os.Args[1])

	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	poll := cfg.Section("main").Key("poll").MustInt(60)
	fmt.Printf("Polltime is %d seconds.\n", poll)

	for t := range time.NewTicker(time.Duration(poll) * time.Second).C {
		if t == t {
		}
		dviData := dvi.GetDviData(cfg)
		influx.SaveToInflux(cfg, dviData)
	}
}
