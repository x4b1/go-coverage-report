package fixtures_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/x4b1/go-coverage-report/pkg/fixtures"
)

func TestIsEven(t *testing.T) {
	for _, tc := range []struct {
		number int
		expect bool
	}{
		{number: 14, expect: true},
		{number: 5, expect: false},
	} {
		tc := tc
		t.Run(fmt.Sprintf("%d even", tc.number), func(t *testing.T) {
			t.Parallel()
			even, err := fixtures.IsEven(tc.number)
			require.NoError(t, err)
			require.Equal(t, tc.expect, even)
		})
	}
}
