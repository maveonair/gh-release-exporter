# gh-release-exporter

This program checks on GitHub if a newer version of a configured program exists (see `config.example.toml`) and exposes the result as a Prometheus Metrics `gh_release_probe_success`.

## Endpoint

- Prometheus Metrics Endpoint: `:9054/metrics` (listening address can be changed through config file, see [Configuration](#configuration)

## Configuration

Create a configuration file like the following:

```toml
listening_addr = "0.0.0.0:9054"

[releases]
[releases.redis]
last_known_version = "7.0.4"
github_repo = "redis/redis"

[releases.strongswan]
last_known_version = "5.9.1"
github_repo = "strongswan/strongswan"
```

## Run

```sh
$ gh-release-exporter -config config.toml
```
