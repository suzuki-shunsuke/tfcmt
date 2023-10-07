package cli

import (
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
	"github.com/urfave/cli/v2"
)

func newConfig(ctx *cli.Context) (config.Config, error) {
	cfg := config.Config{}
	confPath, err := cfg.Find(ctx.String("config"))
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
