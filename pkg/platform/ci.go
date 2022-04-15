package platform

import (
	"fmt"
	"os"
	"strconv"

	"github.com/suzuki-shunsuke/go-ci-env/cienv"
	"github.com/suzuki-shunsuke/tfcmt/pkg/config"
)

func Complement(cfg *config.Config) error {
	if err := complementWithCIEnv(cfg.CI); err != nil {
		return err
	}

	if err := complementCIInfo(cfg.CI); err != nil {
		return err
	}

	return complementWithGeneric(cfg)
}

func complementCIInfo(ci *config.CI) error {
	if ci.PRNumber <= 0 {
		// support suzuki-shunsuke/ci-info
		if prS := os.Getenv("CI_INFO_PR_NUMBER"); prS != "" {
			a, err := strconv.Atoi(prS)
			if err != nil {
				return fmt.Errorf("parse CI_INFO_PR_NUMBER %s: %w", prS, err)
			}
			ci.PRNumber = a
		}
	}
	return nil
}

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

func complementWithCIEnv(ci *config.CI) error {
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

		if ci.PRNumber <= 0 {
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

func complementWithGeneric(cfg *config.Config) error {
	gen := generic{
		param: Param{
			RepoOwner: cfg.Complement.Owner,
			RepoName:  cfg.Complement.Repo,
			SHA:       cfg.Complement.SHA,
			PRNumber:  cfg.Complement.PR,
			Link:      cfg.Complement.Link,
			Vars:      cfg.Complement.Vars,
		},
	}

	if cfg.CI.Owner == "" {
		cfg.CI.Owner = gen.RepoOwner()
	}

	if cfg.CI.Repo == "" {
		cfg.CI.Repo = gen.RepoName()
	}

	if cfg.CI.SHA == "" {
		cfg.CI.SHA = gen.SHA()
	}

	if cfg.CI.PRNumber <= 0 {
		n, err := gen.PRNumber()
		if err != nil {
			return err
		}
		cfg.CI.PRNumber = n
	}

	if cfg.CI.Link == "" {
		cfg.CI.Link = gen.Link()
	}

	vars := gen.Vars()
	for k, v := range cfg.Vars {
		vars[k] = v
	}
	cfg.Vars = vars

	return nil
}
