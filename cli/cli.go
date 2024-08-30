package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	cli "github.com/urfave/cli/v2"

	"github.com/bookandmusic/envsetup/app"
	"github.com/bookandmusic/envsetup/config"
	"github.com/bookandmusic/envsetup/utils"
)

var (
	commonFlags  = []cli.Flag{helpFlag}
	installFlags = []cli.Flag{helpFlag, tagFlag, forceFlag, httpsProxyFlag, githubProxyFlag}
	updateFlags  = []cli.Flag{helpFlag, httpsProxyFlag, githubProxyFlag}
	deleteFlags  = []cli.Flag{helpFlag}
)

// generateSubcommands creates subcommands for a given action
func generateSubcommands(action func(app.Manager, *app.GlobalFlags) error, apps []app.Manager, actionName string, flags []cli.Flag) []*cli.Command {
	var commands []*cli.Command

	for _, mgr := range apps {
		mgr := mgr // capture the loop variable
		commands = append(commands, &cli.Command{
			Name:     mgr.GetName(),
			Usage:    fmt.Sprintf("%s%s", actionName, mgr.GetName()),
			Flags:    flags,
			HideHelp: true,
			Action: func(c *cli.Context) error {
				return action(mgr, &app.GlobalFlags{
					Force:       c.Bool("force"),
					Tag:         c.String("tag"),
					HttpProxy:   c.String("https-proxy"),
					GithubProxy: c.String("github-proxy"),
				})
			},
		})
	}
	return commands
}

// CreateApp initializes the CLI app with commands
func CreateApp() *cli.App {
	// Initialize global configuration
	config.InitConfig()

	apps := []app.Manager{
		app.NewChsrcManager(),
		app.NewVMRManager(),
		app.NewVimrcManager(),
		app.NewOhMyZshManager(),
	}

	commands := []*cli.Command{
		{
			Name:     "list",
			Usage:    "显示所有应用程序",
			Aliases:  []string{"ls"},
			HideHelp: true,
			Flags:    commonFlags,
			Action: func(c *cli.Context) error {
				cfg := utils.TableConfig{
					Header: table.Row{"名称", "描述"},
					Data: []table.Row{
						{"vimrc", "Vim编辑器的配置文件,用于定制编辑器的行为和外观"},
						{"ohmyzsh", "一个增强Zsh配置的开源框架,提供丰富的插件、主题和配置选项"},
						{"vmr", "一个简单、跨平台的版本管理器,用于管理多种 SDK 及其他工具"},
						{"chsrc", "一个全平台的命令行换源工具"},
					},
				}
				utils.RenderTable(&cfg, os.Stdout)
				return nil
			},
		},
		{
			Name:     "install",
			Usage:    "安装应用程序",
			Aliases:  []string{"i"},
			HideHelp: true,
			Flags:    commonFlags,
			Subcommands: generateSubcommands(
				func(mgr app.Manager, flags *app.GlobalFlags) error { return mgr.Install(flags) },
				apps,
				"安装",
				installFlags,
			),
		},
		{
			Name:     "update",
			Usage:    "更新应用程序",
			Aliases:  []string{"u"},
			HideHelp: true,
			Flags:    commonFlags,
			Subcommands: generateSubcommands(
				func(mgr app.Manager, flags *app.GlobalFlags) error { return mgr.Update(flags) },
				apps,
				"更新",
				updateFlags,
			),
		},
		{
			Name:     "delete",
			Usage:    "删除应用程序",
			Aliases:  []string{"d"},
			HideHelp: true,
			Flags:    commonFlags,
			Subcommands: generateSubcommands(
				func(mgr app.Manager, flags *app.GlobalFlags) error { return mgr.Delete(flags) },
				apps,
				"删除",
				deleteFlags,
			),
		},
	}

	return &cli.App{
		Name:     "envsetup",
		Usage:    "配置基本开发环境",
		HideHelp: true,
		Flags:    commonFlags,
		Commands: commands,
	}
}
