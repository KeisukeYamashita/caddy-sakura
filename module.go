package xxx

import (
	cloudns "github.com/anxuanzi/libdns-cloudns"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)

type Provider struct{ *cloudns.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.cloudns",
		New: func() caddy.Module { return &Provider{new(cloudns.Provider)} },
	}
}

func (p *Provider) Provision(ctx caddy.Context) error {
	replacer := caddy.NewReplacer()
	p.Provider.AuthId = replacer.ReplaceAll(p.Provider.AuthId, "")
	p.Provider.SubAuthId = replacer.ReplaceAll(p.Provider.SubAuthId, "")
	p.Provider.AuthPassword = replacer.ReplaceAll(p.Provider.AuthPassword, "")
	return nil
}

func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "auth_id":
				if d.NextArg() {
					p.Provider.AuthId = d.Val()
				} else {
					return d.ArgErr()
				}
			case "sub_auth_id":
				if d.NextArg() {
					p.Provider.SubAuthId = d.Val()
				} else {
					return d.ArgErr()
				}
			case "auth_password":
				if d.NextArg() {
					p.Provider.AuthPassword = d.Val()
				} else {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.AuthId == "" && p.Provider.SubAuthId == "" {
		return d.Err("missing auth id or sub auth id")
	}
	if p.Provider.AuthPassword == "" {
		return d.Err("missing auth password")
	}
	return nil
}
