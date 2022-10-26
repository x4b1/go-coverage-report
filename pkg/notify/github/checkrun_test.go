package github_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/stretchr/testify/require"

	notifier "github.com/xabi93/go-coverage-report/pkg/notify/github"
)

const (
	testOwner = "test-owner"
	testRepo  = "test-repo"

	body = "Coverage: 66.67%"
	sha  = "12314123"
)

func TestCheckRunNotify(t *testing.T) {
	t.Skip("skipping for testing porpouses")

	cliError := errors.New("error")

	customReportName := "my-report"

	for _, tc := range []struct {
		name        string
		checkName   string
		setup       func(*testing.T, *CheckCreatorMock)
		expectedErr error
		assertCalls func(*testing.T, *CheckCreatorMock)
	}{
		{
			name: "fails creating check",
			setup: func(t *testing.T, ccm *CheckCreatorMock) {
				t.Helper()
				ccm.CreateCheckRunFunc = func(context.Context, string, string, github.CreateCheckRunOptions) (*github.CheckRun, *github.Response, error) {
					return nil, nil, cliError
				}
			},
			expectedErr: cliError,
		},
		{
			name: "success",
			setup: func(t *testing.T, ccm *CheckCreatorMock) {
				t.Helper()
				ccm.CreateCheckRunFunc = func(context.Context, string, string, github.CreateCheckRunOptions) (*github.CheckRun, *github.Response, error) {
					return nil, nil, nil
				}
			},
			assertCalls: func(t *testing.T, ccm *CheckCreatorMock) {
				t.Helper()

				require.Len(t, ccm.CreateCheckRunCalls(), 1)
				call := ccm.CreateCheckRunCalls()[0]
				require.Equal(t, testOwner, call.Owner)
				require.Equal(t, testRepo, call.Repo)
				require.Equal(t, github.CreateCheckRunOptions{
					Name:       notifier.DefaultCheckRunName,
					HeadSHA:    sha,
					Status:     github.String("completed"),
					Conclusion: github.String("success"),
					Output: &github.CheckRunOutput{
						Title:   github.String(notifier.DefaultCheckRunName),
						Summary: github.String(body),
					},
				}, call.Opts)
			},
		},
		{
			name:      "success with custom name",
			checkName: customReportName,
			setup: func(t *testing.T, ccm *CheckCreatorMock) {
				t.Helper()
				ccm.CreateCheckRunFunc = func(context.Context, string, string, github.CreateCheckRunOptions) (*github.CheckRun, *github.Response, error) {
					return nil, nil, nil
				}
			},
			assertCalls: func(t *testing.T, ccm *CheckCreatorMock) {
				t.Helper()

				require.Len(t, ccm.CreateCheckRunCalls(), 1)
				call := ccm.CreateCheckRunCalls()[0]
				require.Equal(t, testOwner, call.Owner)
				require.Equal(t, testRepo, call.Repo)
				require.Equal(t, github.CreateCheckRunOptions{
					Name:       customReportName,
					HeadSHA:    sha,
					Status:     github.String("completed"),
					Conclusion: github.String("success"),
					Output: &github.CheckRunOutput{
						Title:   github.String(customReportName),
						Summary: github.String(body),
					},
				}, call.Opts)
			},
		},
	} {

		ccm := &CheckCreatorMock{}

		if tc.setup != nil {
			tc.setup(t, ccm)
		}

		err := notifier.NewCheckRun(nil, ccm, tc.checkName).Notify(context.Background(), nil, body)
		require.ErrorIs(t, err, tc.expectedErr)

		if tc.assertCalls != nil {
			tc.assertCalls(t, ccm)
		}
	}
}
