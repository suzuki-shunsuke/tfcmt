package main

import (
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/cli"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

var version = ""

func main() {
	urfave.Main("tfcmt", version, cli.Run)
}
