package platform

import (
	"fmt"
	"strconv"

	"github.com/suzuki-shunsuke/tfcmt/pkg/domain"
)

type Param struct {
	RepoOwner []domain.ComplementEntry
	RepoName  []domain.ComplementEntry
	SHA       []domain.ComplementEntry
	PRNumber  []domain.ComplementEntry
	Link      []domain.ComplementEntry
}

type generic struct {
	param Param
}

func (gen *generic) render(entries []domain.ComplementEntry) (string, error) {
	var e error
	for _, entry := range entries {
		a, err := entry.Entry()
		if err != nil {
			e = err
			continue
		}
		if a != "" {
			return a, nil
		}
	}
	return "", e
}

func (gen *generic) returnString(entries []domain.ComplementEntry) string {
	s, err := gen.render(entries)
	if err != nil {
		return ""
	}
	return s
}

func (gen *generic) RepoOwner() string {
	return gen.returnString(gen.param.RepoOwner)
}

func (gen *generic) RepoName() string {
	return gen.returnString(gen.param.RepoName)
}

func (gen *generic) SHA() string {
	return gen.returnString(gen.param.SHA)
}

func (gen *generic) Link() string {
	return gen.returnString(gen.param.Link)
}

func (gen *generic) IsPR() bool {
	return gen.returnString(gen.param.PRNumber) != ""
}

func (gen *generic) PRNumber() (int, error) {
	s, err := gen.render(gen.param.PRNumber)
	if err != nil {
		return 0, err
	}
	if s == "" {
		return 0, nil
	}
	b, err := strconv.Atoi(s)
	if err == nil {
		return b, nil
	}
	return 0, fmt.Errorf("parse pull request number as int: %w", err)
}
