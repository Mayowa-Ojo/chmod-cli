package cmd

import (
	"github.com/Mayowa-Ojo/chmod-cli/internal/ui"
	"github.com/urfave/cli/v2"
)

// Execute serves as the cli entry point
func Execute() *cli.App {
	app := &cli.App{
		Name:     "chmod-cli",
		Usage:    "generate file permissions with the bat of an eye",
		Commands: []*cli.Command{},
		Action: func(c *cli.Context) error {

			if err := ui.InitScreen(); err != nil {
				return err
			}

			return nil
		},
	}

	return app
}
