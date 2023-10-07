package github

import (
	"context"
	"testing"
)

func TestPRNumber(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		prNumber int
		ok       bool
		revision string
	}{
		{
			prNumber: 1,
			ok:       true,
			revision: "xxx",
		},
	}

	for _, testCase := range testCases {
		cfg := newFakeConfig()
		client, err := NewClient(context.Background(), &cfg)
		if err != nil {
			t.Fatal(err)
		}
		api := newFakeAPI()
		client.API = &api
		prNumber, err := client.Commits.PRNumber(context.Background(), testCase.revision)
		if (err == nil) != testCase.ok {
			t.Errorf("got error %q", err)
		}
		if prNumber != testCase.prNumber {
			t.Errorf("got %d but want %d", prNumber, testCase.prNumber)
		}
	}
}
