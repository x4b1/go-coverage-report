package fixtures

import (
	"errors"
)

//go:generate go test -count 1 ./... -coverprofile=cover.out

func IsEven(i int) (bool, error) {
	if i == 0 {
		return false, errors.New("number cannot be 0")
	}

	return i%2 == 0, nil
}
