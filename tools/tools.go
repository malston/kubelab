//go:build tools

package tools

// https://go.dev/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md

//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install mvdan.cc/gofumpt
//go:generate go install github.com/daixiang0/gci
//go:generate go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt
//go:generate go install golang.org/x/tools/cmd/goimports
//go:generate go install golang.org/x/lint/golint
//go:generate go install github.com/go-critic/go-critic/cmd/gocritic

// nolint
import (
	// gci
	_ "github.com/daixiang0/gci"
	// gocritic
	_ "github.com/go-critic/go-critic/cmd/gocritic"
	// golangci-lint
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	// gotestfmt
	_ "github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt"
	// counterfeiter
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	// golint
	_ "golang.org/x/lint/golint"
	// goimports
	_ "golang.org/x/tools/cmd/goimports"
	// gofumpt
	_ "mvdan.cc/gofumpt"
)
