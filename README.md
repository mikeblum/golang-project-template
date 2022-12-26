# golang-project-template

![golang-helloworld](https://user-images.githubusercontent.com/3905463/209570840-6b4c3df0-aca4-4de0-899d-ebc823ae0366.png)

A batteries-included Golang project template derived from bootstrapping many Golang projects

## GOTO libraries 📚

`pre-commit`: https://pre-commit.com/

> A framework for managing and maintaining multi-language pre-commit hooks.

`golanglint-ci`: https://golangci-lint.run/

> Fast linters Runner for Go.

`logrus`: https://github.com/sirupsen/logrus

> Structured, pluggable logging for Go.

`viper`: https://github.com/spf13/viper

> Go configuration with fangs.

`testify`: https://github.com/stretchr/testify

A toolkit with common assertions and mocks that plays nicely with the standard library.

## Inspired By 💡

https://github.com/golang-standards/project-layout

🍴 from https://github.com/cloudflare/cloudflare-go/blob/master/.golintci.yaml

🍴 from https://github.com/github/gitignore/blob/main/Go.gitignore

## Run It 🏃

`go run main.go`

## Test It 🧪

Test for coverage and race conditions

`go test -race -covermode=atomic .
/...`

## Lint It 👕

`pre-commit run --all-files --show-diff-on-failure`

## Fork It 🍴

This is a template project for starting your next Golang proj:

https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template

## How To

#### `fmt.Println` is banned in favor of `logrus`

example `fmt.Println("Hello, world")` will throw an error running `golangci-lint run ./...` or `pre-commit`

#### Environment Variables

| Name          | Description   | Default       |
| ------------- | ------------- | ------------- |
| `CONFIG_PATH`   | config.env directory | ./config.env |
| `LOG_LEVEL`  | logging levels: `trace`,`debug`,`info`,`warn` see [ParseLevel(lvl string)](https://github.com/sirupsen/logrus/blob/fdf1618bf7436ec3ee65753a6e2999c335e97221/logrus.go#L25) | `INFO` |
| `LOG_FORMAT` | logging format: `json` or defaults to plaintext | `PLAIN` |
