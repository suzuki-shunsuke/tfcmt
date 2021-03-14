package platform

import (
	"fmt"
	"os"
)

func Get(ciname string) string {
	switch ciname {
	case "circleci", "circle-ci":
		return circleci()
	case "codebuild":
		return codebuild()
	case "github-actions":
		return githubActions()
	case "cloud-build", "cloudbuild":
		return cloudbuild()
	}
	return ""
}

func circleci() string {
	return os.Getenv("CIRCLE_BUILD_URL")
}

func codebuild() string {
	return os.Getenv("CODEBUILD_BUILD_URL")
}

func githubActions() string {
	return fmt.Sprintf(
		"https://github.com/%s/actions/runs/%s",
		os.Getenv("GITHUB_REPOSITORY"),
		os.Getenv("GITHUB_RUN_ID"),
	)
}

func cloudbuild() string {
	return fmt.Sprintf(
		"https://console.cloud.google.com/cloud-build/builds/%s?project=%s",
		os.Getenv("BUILD_ID"),
		os.Getenv("PROJECT_ID"),
	)
}
