package fixtures

import (
	_ "embed"
)

// CoverFile exports a file coverage for testing purposes.
//

var (
	//go:embed cover.out
	CoverFile []byte
	//go:embed default_md_result.md
	DefaultMDResult string
)
