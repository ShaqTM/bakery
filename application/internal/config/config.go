package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server_port string `yaml:"server_port"`
	Server_addr string `yaml:"server_addr"`
}

func ReadConfig() *Config {
	fmt.Println(os.Executable())
	fContent, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
		return &Config{}
	}
	config := Config{}
	err = yaml.Unmarshal(fContent, &config)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
		return &Config{}
	}

	return &config
}
