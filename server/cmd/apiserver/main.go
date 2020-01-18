package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/nickmurr/go-http-rest-api/internal/app/apiserver"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatalf("Error while reading config: %v", err)
	}

	s := apiserver.New(config)

	if err := s.Start(); err != nil {
		log.Fatalf("Error: %v", err)
	}

}
