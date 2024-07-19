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
        "name": "caddy",
        "api_secret": "xxx", // optional
        "api_token": "xxx"   // optional
      }
    }
  }
}
```

Or either you can use the Caddyfile:

```Caddyfile
tls {
  dns sakura {
    api_token  "xxx" // optional
    api_secret "xxx" // optional
  }
}
```

> [!NOTE]
> If you don't provide `api_token` and `api_secret`, the module will try to read them from the environment variables `SAKURACLOUD_ACCESS_TOKEN` and `SAKURACLOUD_ACCESS_TOKEN_SECRET`.
>
> Note that the environment variables have low priority than the Caddyfile configuration.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
