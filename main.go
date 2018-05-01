package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	ini "gopkg.in/ini.v1"
)

var debug = true

// DVILogin is a type for defining the login at DVI Energi webservice
type DVILogin struct {
	Usermail     string `json:"usermail"`
	Userpassword string `json:"userpassword"`
	Fabnr        int    `json:"fabnr"`
	Get          DVIGet `json:"get"`
}

// DVIGet is a type for defining what information to request from DVI Energi webservice
type DVIGet struct {
	Bestgreen int `json:"bestgreen"`
	Sensor    int `json:"sensor"`
	Relay     int `json:"relay"`
	Timer     int `json:"timer"`
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	data := DVILogin{
		Usermail:     cfg.Section("login").Key("usermail").String(),
		Userpassword: cfg.Section("login").Key("userpassword").String(),
		Fabnr:        cfg.Section("login").Key("fabnr").MustInt(),
		Get: DVIGet{
			Bestgreen: cfg.Section("get").Key("bestgreen").MustInt(0),
			Sensor:    cfg.Section("get").Key("sensor").MustInt(0),
			Relay:     cfg.Section("get").Key("relay").MustInt(0),
			Timer:     cfg.Section("get").Key("timer").MustInt(0)}}

	json, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Could not convert data to json: %s\n", err)
	} else {
		if debug {
			// FIXME: mask password
			fmt.Println(jsonPrettyPrint(string(json)))
		}
	}

	response, err := http.Post("https://ws.dvienergi.com/API/", "application/json", bytes.NewBuffer(json))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		if debug {
			// FIXME: mask password if present (failed login)
			fmt.Println(jsonPrettyPrint(string(data)))
		}
	}
}
