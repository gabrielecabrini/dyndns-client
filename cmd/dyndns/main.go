package main

import (
	"dyndns-client/internal/config"
	"dyndns-client/internal/providers"
	"dyndns-client/internal/runner"
	"flag"
	"log"
)

func main() {
	cfgPath := flag.String("config", "./config.yaml", "Configuration file path")
	flag.Parse()
	log.Printf("Loading config: %s\n", *cfgPath)

	cfg, err := config.GetConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	providers, err := providers.GetProviders(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d providers", len(providers))

	runner.Run(*cfg, providers)

}
