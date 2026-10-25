# Security Policy

## Supported Versions

Only the latest release receives security fixes.

## Reporting a Vulnerability

Please do not open a public issue for security problems. Use GitHub's
[private vulnerability reporting](https://github.com/huohen92/anifeed/security/advisories/new)
or contact the maintainer privately. Include reproduction steps, affected
version, and whether the issue requires an authenticated session.

## Deployment Hardening

AniFeed controls a 115 account and can submit, move, rename, and delete cloud
files. Treat it as an administrative service.

- Always enable Web UI login and use a long, unique password.
- Put the service behind HTTPS (a reverse proxy is recommended).
- Prefer binding to `127.0.0.1` unless the service is protected by a trusted
  reverse proxy or VPN.
- Set `server.access_token` when exposing the API outside a trusted network.
- Never commit `config.toml`, `.cookies`, `db.sqlite`, `posters/`, backups, or
  runtime logs. They are ignored by default.
- Restrict file permissions for the configuration directory and Cookies file.
- Review `[[rss]]` URLs and custom HTTP headers before enabling notifications
  or exposing the service to other users.

## Secrets Handled by AniFeed

AniFeed stores or processes:

- 115 Cookies
- Web UI password hash and session tokens
- TMDB API key
- notification webhook URLs
- RSS URLs and optional request headers

The configuration API returns these values to an authenticated administrator.
Do not share the Web UI session, configuration export, or logs publicly.
