# dl

A diesel train that rolls across your terminal — inspired by `sl`.

```
                __@______________________/___\__/___\____
          _____/ ___   ___   ___   ___   ___   ___   ___ |
         / [==] |\o/| |   | |   | |   | |   | |   | |   ||
        |  [**] |===| | @ | | @ | | @ | | @ | | @ | | @ ||
        |  [==] |___| |___| |___| |___| |___| |___| |___||
         \_______________________________________________/
          (oo)(oo)         (oo)(oo)         (oo)(oo)
```

## Install

```sh
brew install PandasWhoCode/tools/dl
```

## Usage

```sh
dl          # normal speed
dl -f       # fast
dl -s       # slow
```

Press `Ctrl+C` to stop early.

## Build from source

Requires [Go 1.22+](https://go.dev/dl/) and [Task](https://taskfile.dev).

```sh
task build   # produces ./dl
task run     # go run .
task tidy    # go mod tidy
```

### Docker

```sh
task docker:build        # builds distroless image tagged dl
docker run -it --rm dl
```

## License

Apache 2.0 — see [LICENSE](LICENSE).
