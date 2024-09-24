# kubelab

## Overview

`kubelab` is a tool designed to simplify and enhance your Kubernetes workflows. It is built using Go and packaged in a lightweight Alpine-based Docker image.

## Prerequisites

- Docker
- Go 1.19 or later
- Git

## Installation

### Docker

To build the Docker image, run:

```sh
docker build -t kubelab .
```

### Go

To build the binary using Go, run:

```sh
go build -o kubelab
```

## Usage

### Docker

To run `kubelab` using Docker:

```sh
docker run --rm kubelab
```

### Go

To run `kubelab` directly:

```sh
./kubelab
```

## CI/CD

This project uses GitHub Actions for continuous integration and delivery. The pipeline includes the following stages:

- **lint**: Lints the code using `golangci-lint`.
- **build**: Builds the Go binary.
- **test**: Runs tests with race detection and code coverage.
- **release**: Releases the project using `goreleaser`.

## Tools

The project uses several Go tools for development, which are tracked in `tools/tools.go`:

- `golangci-lint`
- `gofumpt`
- `gci`
- `gotestfmt`
- `goimports`
- `golint`
- `gocritic`
- `counterfeiter`

## License

This project is licensed under the MIT License.
