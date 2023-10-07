package github

import (
	"context"
	"testing"
)

func TestCommentPost(t *testing.T) {
	t.Parallel()
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
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			client, err := NewClient(context.Background(), &testCase.config)
			if err != nil {
				t.Fatal(err)
			}
			api := newFakeAPI()
			client.API = &api
			err = client.Comment.Post(context.Background(), testCase.body, &testCase.opt)
			if (err == nil) != testCase.ok {
				t.Errorf("got error %q", err)
			}
		})
	}
}
