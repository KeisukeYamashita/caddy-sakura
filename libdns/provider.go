package libdns

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/libdns/libdns"
	"golang.org/x/exp/slices"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/dns"
)

var (
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
)

type Provider struct {
	client *dns.Service

	// API key for Sakura cloud API.
	Secret string `json:"secret,omitempty"`

	// API token for Sakura cloud API.
	Token string `json:"token,omitempty"`
}

func (p *Provider) init(_ context.Context) {
	if p.client == nil {
		if p.Secret == "" {
			p.Secret = os.Getenv(iaas.APIAccessSecretEnvKey)
		}

		if p.Token == "" {
			p.Token = os.Getenv(iaas.APIAccessTokenEnvKey)
		}

		p.client = dns.New(iaas.NewClient(p.Token, p.Secret))
		return
	}
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.init(ctx)

	d, err := p.client.ReadWithContext(ctx, &dns.ReadRequest{
		ID: types.ZoneIs1aID,
	})
	if err != nil {
		return nil, err
	}

	rs := make(iaas.DNSRecords, len(d.Records))
	for _, r := range d.Records {
		records = append(records, libdns.Record{
			ID:    fmt.Sprintf("%s_%s_%s", d.GetName(), r.GetName(), r.GetType()),
			Name:  r.GetName(),
			TTL:   time.Duration(r.GetTTL()),
			Type:  r.GetType().String(),
			Value: r.GetRData(),
		})
	}

	d, err = p.client.UpdateWithContext(ctx, &dns.UpdateRequest{
		ID:      types.ZoneIs1aID,
		Records: rs,
	})

	newRecords := make([]libdns.Record, len(d.Records))
	for _, r := range d.Records {
		newRecords = append(newRecords, libdns.Record{
			ID:    fmt.Sprintf("%s_%s_%s", d.GetName(), r.GetName(), r.GetType()),
			Name:  r.GetName(),
			TTL:   time.Duration(r.GetTTL()),
			Type:  r.GetType().String(),
			Value: r.GetRData(),
		})
	}

	return newRecords, err
}

func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.init(ctx)

	d, err := p.client.ReadWithContext(ctx, &dns.ReadRequest{
		ID: types.ZoneIs1aID,
	})
	if err != nil {
		return nil, err
	}

	rs := make([]libdns.Record, len(d.Records))
	for _, r := range d.Records {
		records = append(records, libdns.Record{
			ID:    fmt.Sprintf("%s_%s_%s", d.GetName(), r.GetName(), r.GetType()),
			Name:  r.GetName(),
			TTL:   time.Duration(r.GetTTL()),
			Type:  r.GetType().String(),
			Value: r.GetRData(),
		})
	}

	remaining := slices.DeleteFunc(rs, func(r libdns.Record) bool {
		return true
	})

	newRecords := make(iaas.DNSRecords, len(remaining))
	for i, r := range remaining {
		newRecords[i] = &iaas.DNSRecord{
			Name:  r.Name,
			RData: r.Value,
			Type:  types.EDNSRecordType(r.Type),
			TTL:   int(r.TTL.Seconds()),
		}
	}

	d, err = p.client.UpdateWithContext(ctx, &dns.UpdateRequest{
		ID:      types.ZoneIs1aID,
		Records: newRecords,
	})
	if err != nil {
		return nil, err
	}

	rs = make([]libdns.Record, len(d.Records))
	for _, r := range d.Records {
		records = append(records, libdns.Record{
			ID:    fmt.Sprintf("%s_%s_%s", d.GetName(), r.GetName(), r.GetType()),
			Name:  r.GetName(),
			TTL:   time.Duration(r.GetTTL()),
			Type:  r.GetType().String(),
			Value: r.GetRData(),
		})
	}

	return rs, nil

}

func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.init(ctx)

	d, err := p.client.ReadWithContext(ctx, &dns.ReadRequest{
		ID: types.ZoneIs1aID,
	})
	if err != nil {
		return nil, err
	}

	records := make([]libdns.Record, len(d.Records))
	for _, r := range d.Records {
		records = append(records, libdns.Record{
			ID:    fmt.Sprintf("%s_%s_%s", d.GetName(), r.GetName(), r.GetType()),
			Name:  r.GetName(),
			TTL:   time.Duration(r.GetTTL()),
			Type:  r.GetType().String(),
			Value: r.GetRData(),
		})
	}

	return records, nil
}

func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.init(ctx)

	rs := make(iaas.DNSRecords, len(records))

	for _, record := range records {
		r := &iaas.DNSRecord{
			Name:  record.Name,
			RData: record.Value,
			Type:  types.EDNSRecordType(record.Type),
			TTL:   int(record.TTL.Seconds()),
		}

		rs = append(rs, r)
	}

	_, err := p.client.UpdateWithContext(ctx, &dns.UpdateRequest{
		ID:      types.ZoneIs1aID,
		Records: rs,
	})

	if err != nil {
		return nil, err
	}

	return records, nil
}
