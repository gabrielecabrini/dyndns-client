package providers

import (
	"dyndns-client/internal/config"
	"dyndns-client/internal/providers/cloudflare"
	"fmt"
	"log"
)

type DNSProvider interface {
	DisplayName() string
	UpdateRecord(ip string) error
}

func GetProviders(cfg config.Config) ([]DNSProvider, error) {
	var providers []DNSProvider

	for name, settings := range cfg.Providers {
		switch name {
		case "cloudflare":
			zone, _ := settings["zone"].(string)
			token, _ := settings["token"].(string)
			domain, _ := settings["domain"].(string)

			if zone == "" || token == "" || domain == "" {
				return nil, fmt.Errorf("invalid config for provider %s: zone, token, and domain are required", name)
			}

			providers = append(providers, &cloudflare.CloudflareProvider{
				Zone:   zone,
				Token:  token,
				Domain: domain,
			})
		default:
			log.Printf("unknown provider %s", name)
			continue
		}
	}

	return providers, nil
}
