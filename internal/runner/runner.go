package runner

import (
	"dyndns-client/internal/config"
	"dyndns-client/internal/providers"
	"dyndns-client/internal/updater"
	"log"
	"time"
)

func Run(cfg config.Config, providers []providers.DNSProvider) error {
	if cfg.Daemon {
		for {
			ip, err := cfg.GetIP()
			if err != nil {
				log.Fatal(err)
			}

			err = updater.Update(providers, ip)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Sleeping for %s", cfg.Interval.String())
			time.Sleep(cfg.Interval)
		}
	} else {
		ip, err := cfg.GetIP()
		if err != nil {
			log.Fatal(err)
		}

		err = updater.Update(providers, ip)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
