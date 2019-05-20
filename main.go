package main

import (
	"fmt"
	"os"

	dvi "github.com/anderskvist/DVIEnergiSmartControl/dvi"
	influx "github.com/anderskvist/DVIEnergiSmartControl/influx"

	ini "gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load(os.Args[1])

	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	dviData := dvi.GetDviData(cfg)
	influx.SaveToInflux(cfg, dviData)
}
