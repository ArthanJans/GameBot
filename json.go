package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var cfg map[string]string

var mem map[string]string

func readJSON(variable interface{}, fileName string) error {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err = json.Unmarshal(raw, variable); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func writeJSON(variable interface{}, fileName string) error {
	b, err := json.Marshal(variable)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err = ioutil.WriteFile(fileName, b, os.FileMode(0)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
