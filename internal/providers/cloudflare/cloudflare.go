package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CloudflareProvider struct {
	Zone   string
	Token  string
	Domain string
}

type dnsRecordRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
	Proxied bool   `json:"proxied"`
}

type zoneResponse struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

type dnsRecordsResponse struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"result"`
}

func (c *CloudflareProvider) DisplayName() string {
	return "Cloudflare"
}

func (c *CloudflareProvider) UpdateRecord(ip string) error {
	zoneID, err := c.getZoneID()
	if err != nil {
		return err
	}

	recordID, err := c.getDNSRecordID(zoneID)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)
	reqBody := dnsRecordRequest{
		Name:    c.Domain,
		Type:    "A",
		Content: ip,
		TTL:     60,
		Comment: "Updated by DynDNS client",
		Proxied: false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response from Cloudflare: %s - %s", resp.Status, string(data))
	}

	return nil
}

func (c *CloudflareProvider) getZoneID() (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%s", c.Zone)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get zone ID: %w", err)
	}
	defer resp.Body.Close()

	var z zoneResponse
	if err := json.NewDecoder(resp.Body).Decode(&z); err != nil {
		return "", fmt.Errorf("failed to decode zone response: %w", err)
	}

	if len(z.Result) == 0 {
		return "", fmt.Errorf("zone not found: %s", c.Zone)
	}

	return z.Result[0].ID, nil
}

func (c *CloudflareProvider) getDNSRecordID(zoneID string) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s&type=A", zoneID, c.Domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get DNS record ID: %w", err)
	}
	defer resp.Body.Close()

	var r dnsRecordsResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("failed to decode DNS records response: %w", err)
	}

	if len(r.Result) == 0 {
		return "", fmt.Errorf("DNS record not found: %s", c.Domain)
	}

	return r.Result[0].ID, nil
}
