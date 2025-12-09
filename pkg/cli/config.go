package cli

import (
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
)

func newConfig(configPath string) (config.Config, error) {
	cfg := config.Config{}
	confPath, err := cfg.Find(configPath)
	if err != nil {
		return cfg, err
	}
	if confPath != "" {
		if err := cfg.LoadFile(confPath); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}
