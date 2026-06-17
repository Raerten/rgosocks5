# RGoSocks5

A lightweight SOCKS5 proxy in Go (`github.com/things-go/go-socks5`) with custom
DNS resolution, IP/FQDN access-control rules, optional username/password auth,
and an optional HTTP status endpoint.

## Layout

- `main.go` — wires config → server (authenticator, rules, resolver, dial stats)
  and handles signals. `version`/`commit`/`date` vars are injected at build time.
- `config/` — all configuration is **environment-driven** via `caarlos0/env`
  struct tags (`Config` in `config/config.go`). Add new settings as fields with
  `env:"..."` + `envDefault:"..."` tags; don't read env vars elsewhere.
- `rules/` — `ProxyRulesSet.Allow` enforces allow/reject by IP CIDR and FQDN, and
  gates BIND/ASSOCIATE. Reject wins over allow; empty allow lists = allow-all.
- `resolver/` — custom DNS resolver (`miekg/dns`) with optional `go-cache` TTL
  cache and IPv6 preference; falls back to the system resolver when `DNS_HOST` is unset.
- `stat/` — connection/byte counters (`sync/atomic`) and the HTTP status server
  (`GET /status`, optional Bearer auth).
- `slogger/` — adapter bridging go-socks5 logging to `log/slog`.
- `version/` — build/version info helper.

## Conventions

- Structured logging only: `log/slog`. No `fmt.Println` for runtime output.
- Configuration via env tags in `config/config.go` — never scatter `os.Getenv`.
- Each package has a matching `_test.go`; keep tests passing for the package you touch.

## Before committing

Run and keep green:

```
gofmt -w .
go vet ./...
go test -race ./...
```

(These also run automatically via `.claude/settings.json` hooks.)

## Build & release

- Release is **tag-driven**: pushing a git tag triggers `.github/workflows/release.yml`
  → GoReleaser builds multi-arch binaries and publishes Docker images
  (`raerten/rgosocks5`). Use the `/release` skill to cut one.
- CI (`.github/workflows/ci.yml`) runs fmt/vet/test/lint on push to `main` and PRs.
- Docker image is `FROM scratch`; the binary must stay statically linked
  (`CGO_ENABLED=0`, already set in `.goreleaser.yaml`).

## Security-sensitive areas

This is a proxy enforcing a trust boundary. Treat changes to `rules/`,
`resolver/`, auth (`main.go`), and `stat/` auth as security-sensitive — use the
`security-reviewer` agent for them.
