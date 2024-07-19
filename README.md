# Sakura module for Caddy

> [Sakura](https://cloud.sakura.ad.jp/) module contains a DNS provider for [Caddy](https://github.com/caddyserver/caddy).

## How to use

Build your Caddy image containing this module:

```console
xcaddy build
    --with github.com/KeisukeYamashita/caddy-xxx
```

## Configuration

To use this module for ACME DNS challenge, configure the Caddy JSON as below:

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
