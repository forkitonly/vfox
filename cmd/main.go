package main

import (
	"github.com/aooohan/version-fox/sdk"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
	"strings"
)

const Version = "0.0.1"

func main() {
	println("VersionFox " + runtime.GOARCH)

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v", "V"},
		Usage:   "print version",
		Action: func(ctx *cli.Context, b bool) error {
			println(Version)
			return nil
		},
	}

	manager := sdk.NewSdkManager()

	app := &cli.App{}
	app.Name = "VersionFox"
	app.Usage = "VersionFox is a tool for sdk version management"
	app.UsageText = "vf [command] [command options]"
	// TODO copyright
	app.Copyright = "TODO Copyright"
	app.Version = Version
	app.Description = "VersionFox is a tool for sdk version management, which allows you to quickly install and use different versions of targeted sdk via the command line."
	app.Suggest = true
	app.Commands = []*cli.Command{
		{
			Name:   "install",
			Usage:  "install a version of sdk",
			Action: sdkVersionParser(manager.Install),
		},
		{
			Name:   "uninstall",
			Usage:  "uninstall a version of sdk",
			Action: sdkVersionParser(manager.Uninstall),
		},
		{
			Name:   "search",
			Usage:  "search a version of sdk",
			Action: sdkVersionParser(manager.Search),
		},
		{
			Name:   "use",
			Usage:  "use a version of sdk",
			Action: sdkVersionParser(manager.Use),
		},
	}

	_ = app.Run(os.Args)

}

func sdkVersionParser(operation func(arg sdk.Arg) error) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		sdkArg := ctx.Args().First()
		if sdkArg == "" {
			return cli.Exit("sdk version is required", 1)
		}
		argArr := strings.Split(sdkArg, "@")
		argsLen := len(argArr)
		if argsLen > 2 {
			return cli.Exit("sdk version is invalid", 1)
		} else if argsLen == 2 {
			return operation(sdk.Arg{
				Name:    strings.ToLower(argArr[0]),
				Version: argArr[1],
			})
		} else {
			return operation(sdk.Arg{
				Name:    strings.ToLower(argArr[0]),
				Version: "latest",
			})
		}
	}
}
