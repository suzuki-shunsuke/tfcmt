package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCommentPost(t *testing.T) { //nolint:tparallel
	t.Setenv("GITHUB_TOKEN", "xxx")
	testCases := []struct {
		name   string
		config Config
		body   string
		opt    PostOptions
		ok     bool
	}{
		{
			name:   "1",
			config: newFakeConfig(),
			body:   "",
			opt: PostOptions{
				Number:   1,
				Revision: "abcd",
			},
			ok: true,
		},
		{
			name:   "2",
			config: newFakeConfig(),
			body:   "",
			opt: PostOptions{
				Number:   2,
				Revision: "",
			},
			ok: true,
		},
		{
			name:   "3",
			config: newFakeConfig(),
			body:   "",
			opt: PostOptions{
				Number:   0,
				Revision: "",
			},
			ok: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			cfg := testCase.config
			client, err := NewClient(t.Context(), &cfg)
			if err != nil {
				t.Fatal(err)
			}
			api := newFakeAPI()
			client.API = &api
			opt := testCase.opt
			err = client.Comment.Post(t.Context(), testCase.body, &opt)
			if (err == nil) != testCase.ok {
				t.Errorf("got error %q", err)
			}
		})
	}
}

func TestCommentList(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "xxx")
	// GitHub GraphQL API can return null in nodes of a connection.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"repository":{"pullRequest":{"comments":{"nodes":[null,{"databaseId":1,"body":"foo","isMinimized":false},null],"pageInfo":{"endCursor":"","hasNextPage":false}}}}}}`))
	}))
	t.Cleanup(srv.Close)

	cfg := newFakeConfig()
	cfg.GraphQLEndpoint = srv.URL
	client, err := NewClient(t.Context(), &cfg)
	if err != nil {
		t.Fatal(err)
	}
	comments, err := client.Comment.List(t.Context(), "owner", "repo", 1)
	if err != nil {
		t.Fatal(err)
	}
	exp := []*IssueComment{
		{DatabaseID: 1, Body: "foo"},
	}
	if diff := cmp.Diff(exp, comments); diff != "" {
		t.Fatal(diff)
	}
}
