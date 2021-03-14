package platform

import (
	"os"
	"reflect"
	"testing"
)

func TestCI(t *testing.T) { //nolint:paralleltest
	testCases := []struct {
		name      string
		envs      []string
		getCI     func() (CI, error)
		testCases []struct {
			name string
			fn   func()
			ci   CI
			ok   bool
		}
	}{
		{
			name: "circleci",
			envs: []string{
				"CIRCLE_SHA1",
				"CIRCLE_BUILD_URL",
				"CIRCLE_PULL_REQUEST",
				"CI_PULL_REQUEST",
				"CIRCLE_PR_NUMBER",
			},
			getCI: circleci,
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("CIRCLE_SHA1", "abcdefg")
						os.Setenv("CIRCLE_BUILD_URL", "https://circleci.com/gh/owner/repo/1234")
						os.Setenv("CIRCLE_PULL_REQUEST", "")
						os.Setenv("CI_PULL_REQUEST", "")
						os.Setenv("CIRCLE_PR_NUMBER", "")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   0,
						},
						URL: "https://circleci.com/gh/owner/repo/1234",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("CIRCLE_SHA1", "abcdefg")
						os.Setenv("CIRCLE_BUILD_URL", "https://circleci.com/gh/owner/repo/1234")
						os.Setenv("CIRCLE_PULL_REQUEST", "https://github.com/owner/repo/pull/1")
						os.Setenv("CI_PULL_REQUEST", "")
						os.Setenv("CIRCLE_PR_NUMBER", "")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   1,
						},
						URL: "https://circleci.com/gh/owner/repo/1234",
					},
					ok: true,
				},
				{
					name: "case 2",
					fn: func() {
						os.Setenv("CIRCLE_SHA1", "abcdefg")
						os.Setenv("CIRCLE_BUILD_URL", "https://circleci.com/gh/owner/repo/1234")
						os.Setenv("CIRCLE_PULL_REQUEST", "")
						os.Setenv("CI_PULL_REQUEST", "2")
						os.Setenv("CIRCLE_PR_NUMBER", "")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   2,
						},
						URL: "https://circleci.com/gh/owner/repo/1234",
					},
					ok: true,
				},
				{
					name: "case 3",
					fn: func() {
						os.Setenv("CIRCLE_SHA1", "abcdefg")
						os.Setenv("CIRCLE_BUILD_URL", "https://circleci.com/gh/owner/repo/1234")
						os.Setenv("CIRCLE_PULL_REQUEST", "")
						os.Setenv("CI_PULL_REQUEST", "")
						os.Setenv("CIRCLE_PR_NUMBER", "3")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   3,
						},
						URL: "https://circleci.com/gh/owner/repo/1234",
					},
					ok: true,
				},
				{
					name: "case 4",
					fn: func() {
						os.Setenv("CIRCLE_SHA1", "")
						os.Setenv("CIRCLE_BUILD_URL", "https://circleci.com/gh/owner/repo/1234")
						os.Setenv("CIRCLE_PULL_REQUEST", "")
						os.Setenv("CI_PULL_REQUEST", "")
						os.Setenv("CIRCLE_PR_NUMBER", "")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "",
							Number:   0,
						},
						URL: "https://circleci.com/gh/owner/repo/1234",
					},
					ok: true,
				},
				{
					name: "case 5",
					fn: func() {
						os.Setenv("CIRCLE_SHA1", "")
						os.Setenv("CIRCLE_BUILD_URL", "https://circleci.com/gh/owner/repo/1234")
						os.Setenv("CIRCLE_PULL_REQUEST", "abcdefg")
						os.Setenv("CI_PULL_REQUEST", "")
						os.Setenv("CIRCLE_PR_NUMBER", "")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "",
							Number:   0,
						},
						URL: "https://circleci.com/gh/owner/repo/1234",
					},
					ok: false,
				},
			},
		}, {
			name: "codebuild",
			envs: []string{
				"CODEBUILD_RESOLVED_SOURCE_VERSION",
				"CODEBUILD_SOURCE_VERSION",
				"CODEBUILD_BUILD_URL",
			},
			getCI: codebuild,
			// https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref.html
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("CODEBUILD_RESOLVED_SOURCE_VERSION", "abcdefg")
						os.Setenv("CODEBUILD_SOURCE_VERSION", "pr/123")
						os.Setenv("CODEBUILD_BUILD_URL", "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   123,
						},
						URL: "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("CODEBUILD_RESOLVED_SOURCE_VERSION", "abcdefg")
						os.Setenv("CODEBUILD_SOURCE_VERSION", "pr/1")
						os.Setenv("CODEBUILD_BUILD_URL", "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   1,
						},
						URL: "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new",
					},
					ok: true,
				},
				{
					name: "case 2",
					fn: func() {
						os.Setenv("CODEBUILD_RESOLVED_SOURCE_VERSION", "")
						os.Setenv("CODEBUILD_SOURCE_VERSION", "")
						os.Setenv("CODEBUILD_BUILD_URL", "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "",
							Number:   0,
						},
						URL: "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new",
					},
					ok: true,
				},
				{
					name: "case 3",
					fn: func() {
						os.Setenv("CODEBUILD_RESOLVED_SOURCE_VERSION", "")
						os.Setenv("CODEBUILD_SOURCE_VERSION", "pr/abc")
						os.Setenv("CODEBUILD_BUILD_URL", "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "",
							Number:   0,
						},
						URL: "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new",
					},
					ok: false,
				},
				{
					name: "case 4",
					fn: func() {
						os.Setenv("CODEBUILD_RESOLVED_SOURCE_VERSION", "")
						os.Setenv("CODEBUILD_SOURCE_VERSION", "f3008ac30d28ac38ae2533c2b153f00041661f22")
						os.Setenv("CODEBUILD_BUILD_URL", "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "",
							Number:   0,
						},
						URL: "https://ap-northeast-1.console.aws.amazon.com/codebuild/home?region=ap-northeast-1#/builds/test:f2ae4314-c2d6-4db6-83c2-eacbab1517b7/view/new",
					},
					ok: true,
				},
			},
		}, {
			name: "github-actions",
			envs: []string{
				"GITHUB_SHA",
				"GITHUB_REPOSITORY",
				"GITHUB_RUN_ID",
			},
			getCI: func() (CI, error) {
				return githubActions(), nil
			},
			// https://help.github.com/ja/actions/configuring-and-managing-workflows/using-environment-variables
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("GITHUB_SHA", "abcdefg")
						os.Setenv("GITHUB_REPOSITORY", "suzuki-shunsuke/tfcmt")
						os.Setenv("GITHUB_RUN_ID", "12345")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   0,
						},
						URL: "https://github.com/suzuki-shunsuke/tfcmt/actions/runs/12345",
					},
					ok: true,
				},
			},
		}, {
			name: "cloudbuild",
			envs: []string{
				"COMMIT_SHA",
				"BUILD_ID",
				"PROJECT_ID",
				"_PR_NUMBER",
			},
			getCI: cloudbuild,
			// https://cloud.google.com/cloud-build/docs/configuring-builds/substitute-variable-values
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("COMMIT_SHA", "abcdefg")
						os.Setenv("BUILD_ID", "build-id")
						os.Setenv("PROJECT_ID", "gcp-project-id")
						os.Setenv("_PR_NUMBER", "123")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   123,
						},
						URL: "https://console.cloud.google.com/cloud-build/builds/build-id?project=gcp-project-id",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("COMMIT_SHA", "")
						os.Setenv("BUILD_ID", "build-id")
						os.Setenv("PROJECT_ID", "gcp-project-id")
						os.Setenv("_PR_NUMBER", "")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "",
							Number:   0,
						},
						URL: "https://console.cloud.google.com/cloud-build/builds/build-id?project=gcp-project-id",
					},
					ok: true,
				},
				{
					name: "case 2",
					fn: func() {
						os.Setenv("COMMIT_SHA", "")
						os.Setenv("BUILD_ID", "build-id")
						os.Setenv("PROJECT_ID", "gcp-project-id")
						os.Setenv("_PR_NUMBER", "abc")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "",
							Number:   0,
						},
						URL: "https://console.cloud.google.com/cloud-build/builds/build-id?project=gcp-project-id",
					},
					ok: false,
				},
			},
		},
	}
	for i, parentTestCase := range testCases { //nolint:paralleltest
		parentTestCase := parentTestCase
		if parentTestCase.name == "" {
			t.Fatalf(`index %d: parentTestCase.name == ""`, i)
		}
		t.Run(parentTestCase.name, func(t *testing.T) {
			saveEnvs := make(map[string]string)
			for _, key := range parentTestCase.envs {
				saveEnvs[key] = os.Getenv(key)
				os.Unsetenv(key)
			}
			defer func() {
				for key, value := range saveEnvs {
					os.Setenv(key, value)
				}
			}()
			for i, testCase := range parentTestCase.testCases {
				testCase := testCase
				if testCase.name == "" {
					t.Fatalf(`index %d: testCase.name == ""`, i)
				}
				t.Run(testCase.name, func(t *testing.T) {
					testCase.fn()
					ci, err := parentTestCase.getCI()
					if !reflect.DeepEqual(ci, testCase.ci) {
						t.Errorf("got %q but want %q", ci, testCase.ci)
					}
					if (err == nil) != testCase.ok {
						t.Errorf("got error %q", err)
					}
				})
			}
		})
	}
}
