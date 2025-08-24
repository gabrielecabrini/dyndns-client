package updater

import (
	"dyndns-client/internal/providers"
	"log"
)

func Update(providers []providers.DNSProvider, ip string) error {
	for _, p := range providers {
		if err := p.UpdateRecord(ip); err != nil {
			return err
		}
		log.Printf("[%s] updated to %s", p.DisplayName(), ip)
	}
	return nil
}
