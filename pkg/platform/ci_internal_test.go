package platform

import (
	"testing"
)

func TestGetLink(t *testing.T) {
	const (
		circleCI         = "circleci"
		codeBuild        = "codebuild"
		gitHubActions    = "github-actions"
		googleCloudBuild = "google-cloud-build"
	)
	testCases := []struct {
		name   string
		ciname string
		env    map[string]string
		exp    string
	}{
		{
			name:   "github actions",
			ciname: gitHubActions,
			env: map[string]string{
				"GITHUB_SERVER_URL":  "https://github.com",
				"GITHUB_REPOSITORY":  "suzuki-shunsuke/tfcmt",
				"GITHUB_RUN_ID":      "1",
				"GITHUB_RUN_ATTEMPT": "2",
			},
			exp: "https://github.com/suzuki-shunsuke/tfcmt/actions/runs/1/attempts/2",
		},
		{
			name:   "github actions without run attempt",
			ciname: gitHubActions,
			env: map[string]string{
				"GITHUB_SERVER_URL":  "https://github.com",
				"GITHUB_REPOSITORY":  "suzuki-shunsuke/tfcmt",
				"GITHUB_RUN_ID":      "1",
				"GITHUB_RUN_ATTEMPT": "",
			},
			exp: "https://github.com/suzuki-shunsuke/tfcmt/actions/runs/1",
		},
		{
			name:   "circle ci",
			ciname: circleCI,
			env: map[string]string{
				"CIRCLE_BUILD_URL": "https://circleci.com/gh/suzuki-shunsuke/tfcmt/1",
			},
			exp: "https://circleci.com/gh/suzuki-shunsuke/tfcmt/1",
		},
		{
			name:   "code build",
			ciname: codeBuild,
			env: map[string]string{
				"CODEBUILD_BUILD_URL": "https://console.aws.amazon.com/codebuild/home",
			},
			exp: "https://console.aws.amazon.com/codebuild/home",
		},
		{
			name:   "google cloud build",
			ciname: googleCloudBuild,
			env: map[string]string{
				"_REGION":    "asia-northeast1",
				"BUILD_ID":   "1",
				"PROJECT_ID": "foo",
			},
			exp: "https://console.cloud.google.com/cloud-build/builds;region=asia-northeast1/1?project=foo",
		},
		{
			name:   "google cloud build without region",
			ciname: googleCloudBuild,
			env: map[string]string{
				"_REGION":    "",
				"BUILD_ID":   "1",
				"PROJECT_ID": "foo",
			},
			exp: "https://console.cloud.google.com/cloud-build/builds;region=global/1?project=foo",
		},
		{
			name:   "unknown ci",
			ciname: "unknown",
			exp:    "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}
			if link := getLink(testCase.ciname); link != testCase.exp {
				t.Fatalf("got %q but want %q", link, testCase.exp)
			}
		})
	}
}
