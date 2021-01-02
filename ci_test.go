package main

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
			name: "travisci",
			envs: []string{
				"TRAVIS_PULL_REQUEST_SHA",
				"TRAVIS_PULL_REQUEST",
				"TRAVIS_COMMIT",
			},
			getCI: travisci,
			// https://docs.travis-ci.com/user/environment-variables/
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("TRAVIS_PULL_REQUEST_SHA", "abcdefg")
						os.Setenv("TRAVIS_PULL_REQUEST", "1")
						os.Setenv("TRAVIS_COMMIT", "hijklmn")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   1,
						},
						URL: "",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("TRAVIS_PULL_REQUEST_SHA", "abcdefg")
						os.Setenv("TRAVIS_PULL_REQUEST", "false")
						os.Setenv("TRAVIS_COMMIT", "hijklmn")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "hijklmn",
							Number:   0,
						},
						URL: "",
					},
					ok: true,
				},
			},
		}, {
			name: "teamcity",
			envs: []string{
				"BUILD_VCS_NUMBER",
				"BUILD_NUMBER",
			},
			getCI: teamcity,
			// https://confluence.jetbrains.com/display/TCD18/Predefined+Build+Parameters
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("BUILD_VCS_NUMBER", "fafef5adb5b9c39244027c8f16f7c3aa7e352b2e")
						os.Setenv("BUILD_NUMBER", "123")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "fafef5adb5b9c39244027c8f16f7c3aa7e352b2e",
							Number:   123,
						},
						URL: "",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("BUILD_VCS_NUMBER", "abcdefg")
						os.Setenv("BUILD_NUMBER", "false")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   0,
						},
						URL: "",
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
			name: "drone",
			envs: []string{
				"DRONE_COMMIT_SHA",
				"DRONE_PULL_REQUEST",
			},
			getCI: drone,
			// https://docs.drone.io/reference/environ/
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("DRONE_COMMIT_SHA", "abcdefg")
						os.Setenv("DRONE_PULL_REQUEST", "1")
						os.Setenv("DRONE_BUILD_LINK", "https://cloud.drone.io/owner/repo/1")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   1,
						},
						URL: "https://cloud.drone.io/owner/repo/1",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("DRONE_COMMIT_SHA", "abcdefg")
						os.Setenv("DRONE_PULL_REQUEST", "")
						os.Setenv("DRONE_BUILD_LINK", "https://cloud.drone.io/owner/repo/1")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   0,
						},
						URL: "https://cloud.drone.io/owner/repo/1",
					},
					ok: true,
				},
				{
					name: "case 2",
					fn: func() {
						os.Setenv("DRONE_COMMIT_SHA", "abcdefg")
						os.Setenv("DRONE_PULL_REQUEST", "abc")
						os.Setenv("DRONE_BUILD_LINK", "https://cloud.drone.io/owner/repo/1")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   0,
						},
						URL: "https://cloud.drone.io/owner/repo/1",
					},
					ok: false,
				},
			},
		}, {
			name: "jenkins",
			envs: []string{
				"GIT_COMMIT",
				"BUILD_URL",
				"PULL_REQUEST_NUMBER",
				"PULL_REQUEST_URL",
			},
			getCI: jenkins,
			// https://wiki.jenkins.io/display/JENKINS/Building+a+software+project#Buildingasoftwareproject-belowJenkinsSetEnvironmentVariables
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("GIT_COMMIT", "abcdefg")
						os.Setenv("PULL_REQUEST_NUMBER", "123")
						os.Setenv("BUILD_URL", "http://jenkins.example.com/jenkins/job/test-job/1")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   123,
						},
						URL: "http://jenkins.example.com/jenkins/job/test-job/1",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("GIT_COMMIT", "abcdefg")
						os.Setenv("PULL_REQUEST_NUMBER", "")
						os.Setenv("PULL_REQUEST_URL", "https://github.com/owner/repo/pull/1111")
						os.Setenv("BUILD_URL", "http://jenkins.example.com/jenkins/job/test-job/123")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   1111,
						},
						URL: "http://jenkins.example.com/jenkins/job/test-job/123",
					},
					ok: true,
				},
				{
					name: "case 2",
					fn: func() {
						os.Setenv("PULL_REQUEST_NUMBER", "")
						os.Setenv("PULL_REQUEST_URL", "")
						os.Setenv("GIT_COMMIT", "abcdefg")
						os.Setenv("BUILD_URL", "http://jenkins.example.com/jenkins/job/test-job/456")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   0,
						},
						URL: "http://jenkins.example.com/jenkins/job/test-job/456",
					},
					ok: true,
				},
			},
		}, {
			name: "jenkins-gitlab",
			envs: []string{
				"BUILD_URL",
				"gitlabBefore",
				"gitlabMergeRequestIid",
			},
			getCI: jenkins,
			// https://wiki.jenkins.io/display/JENKINS/Building+a+software+project#Buildingasoftwareproject-belowJenkinsSetEnvironmentVariables
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("gitlabBefore", "abcdefg")
						os.Setenv("gitlabMergeRequestIid", "123")
						os.Setenv("BUILD_URL", "http://jenkins.example.com/jenkins/job/test-job/1")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   123,
						},
						URL: "http://jenkins.example.com/jenkins/job/test-job/1",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("gitlabMergeRequestIid", "")
						os.Setenv("gitlabBefore", "abcdefg")
						os.Setenv("BUILD_URL", "http://jenkins.example.com/jenkins/job/test-job/456")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   0,
						},
						URL: "http://jenkins.example.com/jenkins/job/test-job/456",
					},
					ok: true,
				},
			},
		}, {
			name: "gitlab",
			envs: []string{
				"CI_COMMIT_SHA",
				"CI_JOB_URL",
				"CI_MERGE_REQUEST_IID",
				"CI_MERGE_REQUEST_REF_PATH",
			},
			getCI: gitlabci,
			// https://docs.gitlab.com/ee/ci/variables/README.html
			testCases: []struct {
				name string
				fn   func()
				ci   CI
				ok   bool
			}{
				{
					name: "case 0",
					fn: func() {
						os.Setenv("CI_COMMIT_SHA", "abcdefg")
						os.Setenv("CI_JOB_URL", "https://gitlab.com/owner/repo/-/jobs/111111111")
						os.Setenv("CI_MERGE_REQUEST_IID", "1")
						os.Setenv("CI_MERGE_REQUEST_REF_PATH", "refs/merge-requests/1/head")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   1,
						},
						URL: "https://gitlab.com/owner/repo/-/jobs/111111111",
					},
					ok: true,
				},
				{
					name: "case 1",
					fn: func() {
						os.Setenv("CI_COMMIT_SHA", "hijklmn")
						os.Setenv("CI_JOB_URL", "https://gitlab.com/owner/repo/-/jobs/222222222")
						os.Setenv("CI_MERGE_REQUEST_REF_PATH", "refs/merge-requests/123/head")
						os.Unsetenv("CI_MERGE_REQUEST_IID")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "hijklmn",
							Number:   123,
						},
						URL: "https://gitlab.com/owner/repo/-/jobs/222222222",
					},
					ok: true,
				},
				{
					name: "case 2",
					fn: func() {
						os.Setenv("CI_COMMIT_SHA", "hijklmn")
						os.Setenv("CI_JOB_URL", "https://gitlab.com/owner/repo/-/jobs/333333333")
						os.Unsetenv("CI_MERGE_REQUEST_IID")
						os.Unsetenv("CI_MERGE_REQUEST_REF_PATH")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "hijklmn",
							Number:   0,
						},
						URL: "https://gitlab.com/owner/repo/-/jobs/333333333",
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
						os.Setenv("GITHUB_REPOSITORY", "mercari/tfnotify")
						os.Setenv("GITHUB_RUN_ID", "12345")
					},
					ci: CI{
						PR: PullRequest{
							Revision: "abcdefg",
							Number:   0,
						},
						URL: "https://github.com/mercari/tfnotify/actions/runs/12345",
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
