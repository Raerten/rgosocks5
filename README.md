Dockerized socks5 proxy written in Go.

Docker image builds for `linux/amd64` `linux/arm/v7` `linux/arm64/v8`

Binary builds for `linux/amd64` `linux/arm/v7` `linux/arm64/v8` `windows`

[![Build status](https://github.com/Raerten/rgosocks5/actions/workflows/release.yml/badge.svg)](https://github.com/Raerten/rgosocks5/actions/workflows/release.yml)
![Go version](https://img.shields.io/github/go-mod/go-version/raerten/rgosocks5)

![Docker Pulls](https://img.shields.io/docker/pulls/raerten/rgosocks5)
![Docker Image Size (tag)](https://img.shields.io/docker/image-size/raerten/rgosocks5/latest)

## docker-compose.yml:

```yml
version: "3"

services:
  socks5:
    image: raerten/rgosocks5
    ports:
      - 1080:1080 # socks5 port
    environment:
      - PROXY_USER=
      - PROXY_PASSWORD=
      - PROXY_PORT=1080
      # Timezone for accurate log times
      - TZ=Europe/Moscow
```

🆕 Image available as [raerten/rgosocks5](https://hub.docker.com/r/raerten/rgosocks5)

## Env variables

| Environment variable    | Description                                                                                  | Default value           |
|-------------------------|----------------------------------------------------------------------------------------------|-------------------------|
| PROXY_USER              | Username for proxy                                                                           |                         |
| PROXY_PASSWORD          | Password for proxy                                                                           |                         |
| PROXY_HOST              | Host for proxy                                                                               | 0.0.0.0                 |
| PROXY_PORT              | Port for proxy                                                                               | 1080                    |
| PROXY_ADDRESS           | Address for proxy                                                                            | $PROXY_HOST:$PROXY_PORT |
| TZ                      | Timezone for accurate log times                                                              | UTC                     |
| LOG_LEVEL_DEBUG         | Enable debug logs                                                                            | false                   |
| PROXY_ALLOWED_DEST_FQDN | Comma separated white list of dest FQDN                                                      |                         |
| PROXY_REJECT_DEST_FQDN  | Comma separated black list of dest FQDN                                                      |                         |
| PROXY_ALLOWED_IPS       | Comma separated white list of dest IP or CIDR                                                |                         |
| PROXY_REJECT_IPS        | Comma separated black list of dest IP or CIDR                                                |                         |
| PROXY_DISABLE_BIND      | Disable bind                                                                                 | false                   |
| PROXY_DISABLE_ASSOCIATE | Disable associate                                                                            | false                   |
| DNS_HOST                | Host for of custom UDP DNS server<br/>If empty - use system resolve                          |                         |
| DNS_PORT                | Port for custom UDP DNS server                                                               | 53                      |
| DNS_USE_CACHE           | Use program cache for custom DNS server<br/>Respect TTL<br/>Works only for custom DNS server | true                    |
| PREFER_IPV6             | Prefer IPv6 IP when resolve FQDN                                                             | false                   |


## License

[![MIT](https://img.shields.io/github/license/raerten/rgosocks5)](https://github.com/raerten/rgosocks5/blob/master/LICENSE)