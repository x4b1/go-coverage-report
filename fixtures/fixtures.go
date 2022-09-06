package fixtures

import (
	_ "embed"
)

// CoverFile exports a file coverage for testing purposes.
//
//go:embed cover.out
var CoverFile []byte
