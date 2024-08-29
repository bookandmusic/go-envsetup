package cli

import (
	cli "github.com/urfave/cli/v2"
)

var (
	helpFlag = &cli.BoolFlag{
		Name:               "help",
		Aliases:            []string{"h"},
		Usage:              "打印此帮助，或 h, help, -h, --help",
		DisableDefaultText: true,
	}
	forceFlag = &cli.BoolFlag{
		Name:    "force",
		Aliases: []string{"f"},
		Usage:   "强制操作",
	}
	tagFlag = &cli.StringFlag{
		Name:    "tag",
		Aliases: []string{"t"},
		Usage:   "指定安装的版本。默认安装最新版本",
	}
	httpsProxyFlag = &cli.StringFlag{
		Name:    "https-proxy",
		Aliases: []string{"hp"},
		Usage:   "为HTTP请求启用代理。示例: --https-proxy=http://127.0.0.1:7890/",
	}
	githubProxyFlag = &cli.StringFlag{
		Name:    "github-proxy",
		Aliases: []string{"gp"},
		Usage:   "为GitHub请求启用代理。示例: --github-proxy=https://mirror.ghproxy.com/",
	}
)
