# golang-project-template

💯 [![Go Report Card](https://goreportcard.com/badge/github.com/mikeblum/golang-project-template)](https://goreportcard.com/report/github.com/mikeblum/golang-project-template)

![Coverage](https://raw.githubusercontent.com/mikeblum/golang-project-template/main/images/badges/coverage.svg)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A 🔋batteries-included Golang project template derived from bootstrapping many Golang projects

---

![golang-helloworld](https://user-images.githubusercontent.com/3905463/209570840-6b4c3df0-aca4-4de0-899d-ebc823ae0366.png)

## GOTO libraries 📚

`pre-commit`: https://pre-commit.com/

> A framework for managing and maintaining multi-language pre-commit hooks.

`golangci-lint`: https://golangci-lint.run/

> Fast linters Runner for Go.

`staticcheck`: https://staticcheck.dev/

> Advanced Go linter, pinned as a build tool under `tools/` and run with `make staticcheck`.

`log/slog`: 

> Stdlib structured, pluggable logging >= Go 1.21.0.

`koanf`: https://github.com/knadh/koanf

> Simple, lightweight, extensible, configuration management library for Go.

`testify`: https://github.com/stretchr/testify

> A toolkit with common assertions and mocks that plays nicely with the standard library.

## Inspired By 💡

https://github.com/golang-standards/project-layout

🍴 from https://github.com/cloudflare/cloudflare-go/blob/master/.golintci.yaml

🍴 from https://github.com/github/gitignore/blob/main/Go.gitignore

## Run It 🏃

`go run main.go`

## Configure It ☑️

‼️ ⚠️ make sure `.env` is excluded from git as environment variables can contain credentials and secrets.

`vim .env`

```
LOG_LEVEL=debug
LOG_FORMAT=json
```

or via CLI:

`export LOG_FORMAT=json && export LOG_LEVEL=debug && go run main.go`

## Test It 🧪

Test for coverage and race conditions

`go test -race -covermode=atomic .
/...`

## Lint It 👕

`pre-commit run --all-files --show-diff-on-failure`

## Make It ⚙️

`make help`

## Secure It 

`make vuln` 🔒

With [Vulnerability Management for Go](https://go.dev/blog/vuln) scan our proj against the Go vulnerability database.

## Fork It 🍴

This is a template project for starting your next Golang proj:

https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template

## How To

### `fmt.Println` is banned in favor of `log/slog`

example `fmt.Println("Hello, world")` will throw an error running `golangci-lint run ./...` or `pre-commit`

### Environment Variables

**note:** env variable values are case-insensitive ex. `LOG_LEVEL=` both `INFO` and `info` are valid.

| Name          | Description   | Default       |
| ------------- | ------------- | ------------- |
| `CONF_PATH`   |  dotenv-style conf file | `./.env` |
| `LOG_LEVEL`  | logging levels: `trace`,`debug`,`info`,`warn` | `INFO` |
| `LOG_FORMAT` | logging format: `json` or defaults to plaintext | `TEXT` |


## Roadmap

Branch for web-based Golang apps with Gin 🥃

https://github.com/gin-gonic/gin

Branch for cli-based Golang apps with `urfave/cli`

https://github.com/urfave/cli/
