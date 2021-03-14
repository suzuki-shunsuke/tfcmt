package config

import (
	"errors"

	"github.com/suzuki-shunsuke/tfcmt/pkg/domain"
)

type Complement struct {
	PR    []domain.ComplementEntry
	Owner []domain.ComplementEntry
	Repo  []domain.ComplementEntry
	SHA   []domain.ComplementEntry
	Link  []domain.ComplementEntry
}

type rawComplement struct {
	PR    []map[string]interface{}
	Owner []map[string]interface{}
	Repo  []map[string]interface{}
	SHA   []map[string]interface{}
	Link  []map[string]interface{}
}

func convComplementEntries(maps []map[string]interface{}) ([]domain.ComplementEntry, error) {
	entries := make([]domain.ComplementEntry, len(maps))
	for i, m := range maps {
		entry, err := convComplementEntry(m)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}
	return entries, nil
}

func convComplementEntry(m map[string]interface{}) (domain.ComplementEntry, error) {
	t, ok := m["type"]
	if !ok {
		return nil, errors.New(`"type" is required`)
	}
	typ, ok := t.(string)
	if !ok {
		return nil, errors.New(`"type" must be string`)
	}
	switch typ {
	case "envsubst":
		entry := ComplementEnvsubstEntry{}
		if err := newComplementEnvsubstEntry(m, &entry); err != nil {
			return nil, err
		}
		return &entry, nil
	case "template":
		entry := ComplementTemplateEntry{}
		if err := newComplementTemplateEntry(m, &entry); err != nil {
			return nil, err
		}
		return &entry, nil
	default:
		return nil, errors.New(`unsupported type: ` + typ)
	}
}

func (cpl *Complement) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var val rawComplement
	if err := unmarshal(&val); err != nil {
		return err
	}

	pr, err := convComplementEntries(val.PR)
	if err != nil {
		return err
	}
	cpl.PR = pr

	owner, err := convComplementEntries(val.Owner)
	if err != nil {
		return err
	}
	cpl.Owner = owner

	repo, err := convComplementEntries(val.Repo)
	if err != nil {
		return err
	}
	cpl.Repo = repo

	sha, err := convComplementEntries(val.SHA)
	if err != nil {
		return err
	}
	cpl.SHA = sha

	link, err := convComplementEntries(val.Link)
	if err != nil {
		return err
	}
	cpl.Link = link

	return nil
}
