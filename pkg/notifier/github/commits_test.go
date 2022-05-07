package github

import (
	"context"
	"testing"
)

func TestMergedPRNumber(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		prNumber int
		ok       bool
		revision string
	}{
		{
			prNumber: 1,
			ok:       true,
			revision: "Merge pull request #1 from suzuki-shunsuke/tfcmt",
		},
		{
			prNumber: 123,
			ok:       true,
			revision: "Merge pull request #123 from suzuki-shunsuke/tfcmt",
		},
		{
			prNumber: 0,
			ok:       false,
			revision: "destroyed the world",
		},
		{
			prNumber: 0,
			ok:       false,
			revision: "Merge pull request #string from suzuki-shunsuke/tfcmt",
		},
	}

	for _, testCase := range testCases {
		cfg := newFakeConfig()
		client, err := NewClient(context.Background(), cfg)
		if err != nil {
			t.Fatal(err)
		}
		api := newFakeAPI()
		client.API = &api
		prNumber, err := client.Commits.MergedPRNumber(context.Background(), testCase.revision)
		if (err == nil) != testCase.ok {
			t.Errorf("got error %q", err)
		}
		if prNumber != testCase.prNumber {
			t.Errorf("got %q but want %q", prNumber, testCase.prNumber)
		}
	}
}
