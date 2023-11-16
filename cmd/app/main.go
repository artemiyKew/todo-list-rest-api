package main

import (
	"flag"
	"log"
	"os"

	"github.com/artemiyKew/todo-list-rest-api/internal/api"
	"gopkg.in/yaml.v2"
)

var config []byte

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.yml", "path to config file")
	flag.Parse()

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	config = data
}

func main() {
	cfg := api.NewConfig()
	err := yaml.Unmarshal(config, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := api.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
