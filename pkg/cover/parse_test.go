package cover_test

import (
	"context"
	"io/fs"
	"io/ioutil"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xabi93/go-coverage-report/pkg/cover"
	"github.com/xabi93/go-coverage-report/pkg/fixtures"
)

func TestFileParser(t *testing.T) {
	fileName := t.TempDir() + "/cover.out"
	require.NoError(t, ioutil.WriteFile(fileName, fixtures.CoverFile, 0o755))

	t.Run("file not exists", func(t *testing.T) {
		_, err := cover.NewFileParser("/unknown_file.out").Parse(context.TODO())
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("success", func(t *testing.T) {
		report, err := cover.NewFileParser(fileName).Parse(context.TODO())
		require.NoError(t, err)

		require.Equal(t, 3, report.TotalStmts)
		require.Equal(t, 1, report.Uncovered)
		require.EqualValues(t, 67, math.Round(report.Coverage()))

		require.Len(t, report.Files, 1)
		file := report.Files[0]
		require.Equal(t, 3, file.TotalStmts)
		require.Equal(t, 1, file.Uncovered)
		require.EqualValues(t, 67, math.Round(file.Coverage()))

		require.Len(t, file.Lines, 1)
		line := file.Lines[0]
		require.Equal(t, 10, line.Start)
		require.Equal(t, 12, line.End)
	})
}
