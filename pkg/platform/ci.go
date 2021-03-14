package platform

import (
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/go-ci-env/cienv"
	"github.com/suzuki-shunsuke/tfcmt/pkg/config"
)

func getLink(ciname string) string {
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

func Complement(ci *config.CI) error {
	if pt := cienv.Get(); pt != nil {
		ci.Name = pt.CI()

		if ci.Owner == "" {
			ci.Owner = pt.RepoOwner()
		}

		if ci.Repo == "" {
			ci.Repo = pt.RepoName()
		}

		if ci.SHA == "" {
			ci.SHA = pt.SHA()
		}

		if ci.PRNumber == 0 {
			n, err := pt.PRNumber()
			if err != nil {
				return err
			}
			ci.PRNumber = n
		}

		if ci.Link == "" {
			ci.Link = getLink(ci.Name)
		}
	}
	return nil
}
