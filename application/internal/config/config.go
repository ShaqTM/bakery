package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server_port            string `yaml:"SERVER_PORT"`
	Server_addr            string `yaml:"SERVER_ADDR"`
	PG_connect_string      string `yaml:"PG_CONNECT_STRING"`
	PG_connect_string_init string `yaml:"PG_CONNECT_STRING_INIT"`
	Service_name           string `yaml:"SERVICE_NAME"`
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
