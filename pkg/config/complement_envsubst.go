package config

import (
	"errors"

	"github.com/drone/envsubst"
)

type ComplementEnvsubstEntry struct {
	Value string
}

func (entry *ComplementEnvsubstEntry) Type() string {
	return "envsubst"
}

func (entry *ComplementEnvsubstEntry) Entry() (string, error) {
	return envsubst.EvalEnv(entry.Value)
}

func newComplementEnvsubstEntry(m map[string]interface{}, entry *ComplementEnvsubstEntry) error {
	v, ok := m["value"]
	if !ok {
		return errors.New(`"value" is required`)
	}
	val, ok := v.(string)
	if !ok {
		return errors.New(`"value" must be string`)
	}
	entry.Value = val
	return nil
}
