//go:build tools

package tools

// Manage tool dependencies via go.mod.
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/golang/go/issues/25922
//
// nolint
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/matryer/moq"
)
