# Go flagenv

Uses reflection and struct tags to parse flags and environment variables (also from .env file) with
optional defaults and support for required variables that will cause the parser to panic if they are missing or empty

## Usage

```go
// TODO: Add examples here
```

## Development

Issues and pull requests are welcome

**Requires at least Go 1.26.0**

1. Install dependencies
    ```sh
    go mod download
    ```

1. Run tests
    ```sh
    go test -count=1 -failfast -v ./...
    ```
