package main

import (
	"appimage-cli-tool/cmd/commands"
	"appimage-cli-tool/cmd/commands/install"
	"appimage-cli-tool/cmd/commands/update"
	"github.com/alecthomas/kong"
)

var cli struct {
	Debug      bool `help:"Enable debug mode."`
	DesktopInt bool `help:"Enable desktop integration."`

	Search  commands.SearchCmd `cmd help:"Search applications in the store."`
	Install install.InstallCmd `cmd help:"Install an application."`
	List    commands.ListCmd   `cmd help:"List installed applications."`
	Remove  commands.RemoveCmd `cmd help:"Remove an application."`
	Update  update.UpdateCmd   `cmd help:"Update an application."`
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&commands.Context{Debug: cli.Debug, DesktopIntegration: cli.DesktopInt})
	ctx.FatalIfErrorf(err)

}
