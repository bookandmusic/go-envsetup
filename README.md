# EnvSetup

## 项目简介

**EnvSetup** 是一个用 Go 语言编写的工具，专门用于自动化配置新系统的开发环境。该项目旨在帮助开发者快速、简便地完成一系列常见工具和配置文件的安装和设置。

### 特性

- 自动安装和配置 **Vim** 及 **The Ultimate vimrc**，带来更强大的 Vim 使用体验。
- 安装 **Zsh** 并设置 **Oh My Zsh** 主题和插件。
- 安装和配置 **chsrc**，一个全平台的命令行换源工具。
- 安装和配置 **VMR**，一个简单、跨平台的版本管理器，用于管理多种 SDK 及其他工具。

## 使用说明

1. 前往 [Releases](https://github.com/bookandmusic/go-envsetup/releases) 页面，下载适用于你操作系统的最新版本的 `envsetup` 二进制文件。

2. 使用 `install` 命令将文件移动到系统的 `PATH` 目录中，并设置正确的权限：

    ```bash
    sudo install -m 755 envsetup /usr/local/bin/
    ```

3. 运行工具

    ```bash
    envsetup -h
    ```

