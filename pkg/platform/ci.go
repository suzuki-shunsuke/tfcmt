package platform

import (
	"fmt"
	"os"
)

func GetLink(ciname string) string {
	switch ciname {
	case "circleci", "circle-ci":
		return os.Getenv("CIRCLE_BUILD_URL")
	case "codebuild":
		return os.Getenv("CODEBUILD_BUILD_URL")
	case "github-actions":
		return fmt.Sprintf(
			"https://github.com/%s/actions/runs/%s",
			os.Getenv("GITHUB_REPOSITORY"),
			os.Getenv("GITHUB_RUN_ID"),
		)
	case "cloud-build", "cloudbuild":
		return fmt.Sprintf(
			"https://console.cloud.google.com/cloud-build/builds/%s?project=%s",
			os.Getenv("BUILD_ID"),
			os.Getenv("PROJECT_ID"),
		)
	}
	return ""
}
