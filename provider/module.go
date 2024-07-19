package sakura

import (
	sakura "github.com/KeisukeYamashita/caddy-sakura/libdns"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)

type Provider struct{ *sakura.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.sakura",
		New: func() caddy.Module { return &Provider{new(sakura.Provider)} },
	}
}

func (p *Provider) Provision(ctx caddy.Context) error {
	replacer := caddy.NewReplacer()
	p.Provider.Secret = replacer.ReplaceAll(p.Provider.Secret, "")
	p.Provider.Token = replacer.ReplaceAll(p.Provider.Token, "")
	return nil
}

func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_secret":
				if d.NextArg() {
					p.Provider.Secret = d.Val()
				} else {
					return d.ArgErr()
				}
			case "api_token":
				if d.NextArg() {
					p.Provider.Token = d.Val()
				} else {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized sub directive '%s'", d.Val())
			}
		}
	}
	return nil
}
