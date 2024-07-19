module for Caddy

> Module contains a DNS provider for [Caddy](https://github.com/caddyserver/caddy).

## How to use

Build your Caddy image containing this module:

```console
xcaddy build
    --with github.com/KeisukeYamashita/caddy-xxx
```

## Configuration

To use this module for ACME DNS challeng, configure the Caddy JSON as below:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "xxx",
      }
    }
  }
}
```

Or either you can use the Caddyfile:

```Caddyfile
tls {
  dns xxx {
    xxx
  }
}
```

