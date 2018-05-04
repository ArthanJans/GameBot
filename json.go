package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var cfg map[string]string

var mem map[string]string

func readJSON(variable interface{}, fileName string) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(raw, variable)
}

func writeJSON(variable interface{}, fileName string) {
	b, err := json.Marshal(variable)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(fileName, b, os.FileMode(0))
}
