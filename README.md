# EnvSetup

## 项目简介

**EnvSetup** 是一个用 Go 语言编写的工具，专门用于自动化配置新系统的开发环境。该项目旨在帮助开发者快速、简便地完成一系列常见工具和配置文件的安装和设置。

### 特性

- 自动安装和配置 **Vim** 及 **The Ultimate vimrc**，带来更强大的 Vim 使用体验。
- 安装 **Zsh** 并设置 **Oh My Zsh** 主题和插件。
- 安装和配置 **chsrc**，一个全平台的命令行换源工具。
- 安装和配置 **VMR**，一个简单、跨平台的版本管理器，用于管理多种 SDK 及其他工具。

## 使用说明

1. 前往 [Releases](https://github.com/yourusername/EnvSetup/releases) 页面，下载适用于你操作系统的最新版本的 `envsetup` 二进制文件。

2. 使用 `install` 命令将文件移动到系统的 `PATH` 目录中，并设置正确的权限：

    ```bash
    sudo install -m 755 envsetup /usr/local/bin/
    ```

3. 运行工具：

    ```bash
    envsetup
    ```

## 功能模块

- Vim 和 vimrc: 自动安装 Vim 并应用预定义的配置文件。
- Zsh 和 Oh My Zsh: 安装 Zsh，并为其设置 Oh My Zsh 主题和常用插件。
- chsrc: 全平台命令行换源工具。
- VMR: 是一款简单，跨平台，且经过良好设计的版本管理器，用于管理多种SDK以及其他工具。

## 贡献

欢迎社区贡献代码、文档以及提供反馈。如果你有任何改进建议或发现问题，请通过提交 issue 或 pull request 的方式与我们联系。

## 许可证

本项目遵循 MIT 许可证。详细信息请参阅 LICENSE 文件。