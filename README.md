[![status][ci-status-badge]][ci-status]
[![PkgGoDev][pkg-go-dev-badge]][pkg-go-dev]

# waitmysql

waitmysql is a command that awaits MySQL server is alive.

It accepts some retry options.

So you can use this command instead of health check options.

## Synopsis

```sh
go install github.com/aereal/waitmysql/cmd/waitmysql@latest
waitmysql -dsn '...'
```

## License

See LICENSE file.

[pkg-go-dev]: https://pkg.go.dev/github.com/aereal/waitmysql
[pkg-go-dev-badge]: https://pkg.go.dev/badge/aereal/waitmysql
[ci-status-badge]: https://github.com/aereal/waitmysql/workflows/CI/badge.svg?branch=main
[ci-status]: https://github.com/aereal/waitmysql/actions/workflows/CI
