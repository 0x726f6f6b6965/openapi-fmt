# openapi-fmt

`openapi-fmt` is a command-line tool for formatting OpenAPI documents. It supports removing custom extension fields from OpenAPI files.

## Features

- Remove all extensions (fields starting with `x-`) from OpenAPI documents, with the option to keep specific fields.
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

## Usage

```bash
o-fmt remove-exts -i <input-file> -o <output-file> -e <exclude1>,<exclude2>
```

- `-i, --input`: Path to the input OpenAPI file (required)
- `-o, --output`: Path to the output OpenAPI file (required)
- `-e, --excludes`: Extension fields to keep (comma-separated, optional)

### Example

Remove all extensions except `x-keep-me`:

```bash
o-fmt remove-exts -i api.yaml -o api.cleaned.yaml -e x-keep-me
```

## Contributing

PRs and issues are welcome!

## License

MIT License