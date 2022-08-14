package input

import "golang.org/x/tools/cover"

// NewLocal creates a new coverage reader for local file.
func NewLocal(filePath string) *Local {
	return &Local{filePath}
}

// Local is a coverage reader for local file.
type Local struct {
	filePath string
}

// Load parses local coverage file and returns a slice of cover.Profile.
func (l Local) Load() ([]*cover.Profile, error) {
	return cover.ParseProfiles(l.filePath)
}
