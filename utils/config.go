package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config struct
type Config struct {
	APIKey   string `json:"api-key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ReadConfig read a json config file
func ReadConfig(configFilename string) (*Config, error) {
	jsonFile, err := os.Open(configFilename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result Config
	json.Unmarshal([]byte(byteValue), &result)

	return &result, nil
}
