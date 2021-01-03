package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// CI represents a common information obtained from all CI platforms
type CI struct {
	PR  PullRequest
	URL string
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Revision string
	Number   int
}

func circleci() (ci CI, err error) {
	ci.PR.Number = 0
	ci.PR.Revision = os.Getenv("CIRCLE_SHA1")
	ci.URL = os.Getenv("CIRCLE_BUILD_URL")
	pr := os.Getenv("CIRCLE_PULL_REQUEST")
	if pr == "" {
		pr = os.Getenv("CI_PULL_REQUEST")
	}
	if pr == "" {
		pr = os.Getenv("CIRCLE_PR_NUMBER")
	}
	if pr == "" {
		return ci, nil
	}
	re := regexp.MustCompile(`[1-9]\d*$`)
	ci.PR.Number, err = strconv.Atoi(re.FindString(pr))
	if err != nil {
		return ci, fmt.Errorf("%v: cannot get env", pr)
	}
	return ci, nil
}

func codebuild() (ci CI, err error) {
	ci.PR.Number = 0
	ci.PR.Revision = os.Getenv("CODEBUILD_RESOLVED_SOURCE_VERSION")
	ci.URL = os.Getenv("CODEBUILD_BUILD_URL")
	sourceVersion := os.Getenv("CODEBUILD_SOURCE_VERSION")
	if sourceVersion == "" {
		return ci, nil
	}
	if !strings.HasPrefix(sourceVersion, "pr/") {
		return ci, nil
	}
	pr := strings.Replace(sourceVersion, "pr/", "", 1)
	if pr == "" {
		return ci, nil
	}
	ci.PR.Number, err = strconv.Atoi(pr)
	return ci, err
}

func githubActions() (ci CI) {
	ci.URL = fmt.Sprintf(
		"https://github.com/%s/actions/runs/%s",
		os.Getenv("GITHUB_REPOSITORY"),
		os.Getenv("GITHUB_RUN_ID"),
	)
	ci.PR.Revision = os.Getenv("GITHUB_SHA")
	return ci
}

func cloudbuild() (ci CI, err error) {
	ci.PR.Number = 0
	ci.PR.Revision = os.Getenv("COMMIT_SHA")
	ci.URL = fmt.Sprintf(
		"https://console.cloud.google.com/cloud-build/builds/%s?project=%s",
		os.Getenv("BUILD_ID"),
		os.Getenv("PROJECT_ID"),
	)
	pr := os.Getenv("_PR_NUMBER")
	if pr == "" {
		return ci, nil
	}
	ci.PR.Number, err = strconv.Atoi(pr)
	return ci, err
}
