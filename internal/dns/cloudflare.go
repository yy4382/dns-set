package dns

import (
	"context"
	"fmt"
	"net"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareProvider struct {
	api *cloudflare.API
}

func NewCloudflareProvider(apiToken string) (*CloudflareProvider, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudflare API client: %w", err)
	}

	return &CloudflareProvider{api: api}, nil
}

func (c *CloudflareProvider) UpdateRecord(domain string, recordType RecordType, ip net.IP, ttl int) error {
	ctx := context.Background()

	zoneID, err := c.getZoneID(ctx, domain)
	if err != nil {
		return fmt.Errorf("failed to get zone ID for domain %s: %w", domain, err)
	}

	recordName := domain
	records, _, err := c.api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Name: recordName,
		Type: string(recordType),
	})
	if err != nil {
		return fmt.Errorf("failed to list DNS records: %w", err)
	}

	ipStr := ip.String()
	
	if len(records) == 0 {
		_, err = c.api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.CreateDNSRecordParams{
			Name:    recordName,
			Type:    string(recordType),
			Content: ipStr,
			TTL:     ttl,
		})
		if err != nil {
			return fmt.Errorf("failed to create DNS record: %w", err)
		}
		return nil
	}

	for _, record := range records {
		if record.Content == ipStr {
			continue
		}

		_, err = c.api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Name:    recordName,
			Type:    string(recordType),
			Content: ipStr,
			TTL:     ttl,
		})
		if err != nil {
			return fmt.Errorf("failed to update DNS record: %w", err)
		}
	}

	return nil
}

func (c *CloudflareProvider) ListRecords(domain string) ([]Record, error) {
	ctx := context.Background()

	zoneID, err := c.getZoneID(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone ID for domain %s: %w", domain, err)
	}

	cfRecords, _, err := c.api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Name: domain,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list DNS records: %w", err)
	}

	var records []Record
	for _, cfRecord := range cfRecords {
		if cfRecord.Type == "A" || cfRecord.Type == "AAAA" {
			records = append(records, Record{
				ID:      cfRecord.ID,
				Name:    cfRecord.Name,
				Type:    RecordType(cfRecord.Type),
				Content: cfRecord.Content,
				TTL:     cfRecord.TTL,
			})
		}
	}

	return records, nil
}

func (c *CloudflareProvider) getZoneID(ctx context.Context, domain string) (string, error) {
	zones, err := c.api.ListZones(ctx, domain)
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %w", err)
	}

	if len(zones) == 0 {
		return "", fmt.Errorf("no zone found for domain %s", domain)
	}

	return zones[0].ID, nil
}

func (c *CloudflareProvider) Name() string {
	return "Cloudflare"
}