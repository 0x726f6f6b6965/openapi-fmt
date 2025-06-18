# openapi-fmt

[![GoDoc](https://godoc.org/github.com/0x726f6f6b6965/openapi-fmt?status.svg)](https://godoc.org/github.com/0x726f6f6b6965/openapi-fmt)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x726f6f6b6965/openapi-fmt)](https://goreportcard.com/report/github.com/0x726f6f6b6965/openapi-fmt)
[![codecov](https://codecov.io/gh/0x726f6f6b6965/openapi-fmt/branch/main/graph/badge.svg)](https://codecov.io/gh/0x726f6f6b6965/openapi-fmt)

`openapi-fmt` is a command-line tool for formatting OpenAPI documents. It supports removing custom extension fields from OpenAPI files and supports extracting OpenAPI from specific paths.

## Features

- Remove all extensions (fields starting with `x-`) from OpenAPI documents, with the option to keep specific fields.
- Split OpenAPI documents by path.
- Supports OpenAPI 3.0 documents in YAML format.

## Installation

```bash
go install github.com/0x726f6f6b6965/openapi-fmt/cmd/o-fmt@latest
```

Or clone and build manually:

```bash
git clone https://github.com/0x726f6f6b6965/openapi-fmt.git
cd openapi-fmt/cmd/o-fmt
go build -o o-fmt
```

Or using `go tool`:

```bash
go get -tool github.com/0x726f6f6b6965/openapi-fmt/cmd/o-fmt@latest
# this will then modify your `go.mod`
```

From there, you can use `go:generate`

## Usage

```bash
o-fmt -i <input-file> -o <output-file> -e <exclude1>,<exclude2>
```

- `-i, --input`: Path to the input OpenAPI file (required)
- `-o, --output`: Path to the output OpenAPI file (required)
- `-e, --excludes`: Extension fields to keep (comma-separated, optional)

### Example

Remove all extensions except `x-keep-me`:

```bash
o-fmt -i api.yaml -o api.cleaned.yaml -e x-keep-me
```

## Contributing

PRs and issues are welcome!

## License

MIT License