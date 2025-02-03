package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func newHelpAll(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:   "help-all",
		Hidden: true,
		Usage:  "show all help",
		Action: func(ctx *cli.Context) error {
			fmt.Fprintln(app.Writer, "```console")
			fmt.Fprintf(app.Writer, "$ %s --help\n", app.Name)
			if err := cli.ShowAppHelp(ctx); err != nil {
				return err
			}
			fmt.Fprintln(app.Writer, "```")
			subcommands := ctx.Command.Subcommands
			ctx.Command.Subcommands = nil
			defer func() {
				ctx.Command.Subcommands = subcommands
			}()
			ignoredCommands := map[string]struct{}{
				"help":     {},
				"help-all": {},
			}
			for _, cmd := range app.Commands {
				if _, ok := ignoredCommands[cmd.Name]; ok {
					continue
				}
				fmt.Fprintf(app.Writer, "\n## %s %s\n\n", app.Name, cmd.Name)
				fmt.Fprintln(app.Writer, "```console")
				fmt.Fprintf(app.Writer, "$ %s %s --help\n", app.Name, cmd.Name)
				if err := cli.ShowCommandHelp(ctx, cmd.Name); err != nil {
					return err
				}
				fmt.Fprintln(app.Writer, "```")
			}
			return nil
		},
	}
}
