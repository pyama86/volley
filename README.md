# volley

snmptrapをwakerに通知するラッパーです。

## Usage

- snmptrapd
```bash
traphandle default /usr/local/bin/volley -w https://waker.example.com/topics/XXXXX/alertmanager.json
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/pyama86/volley
```

## Contribution

1. Fork ([https://github.com/pyama86/volley/fork](https://github.com/pyama86/volley/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pyama86](https://github.com/pyama86)
