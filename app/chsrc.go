package app

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bookandmusic/envsetup/config"
	"github.com/bookandmusic/envsetup/utils"
)

// Define ChsrcManager to handle chsrc operations
type ChsrcManager struct {
	Name     string
	config   *config.Config
	chsrcDir string
}

func NewChsrcManager() *ChsrcManager {
	config := config.GetConfig()

	chsrcDir := fmt.Sprintf("%s/.local/bin", config.HomeDir)
	return &ChsrcManager{
		Name:     "chsrc",
		config:   config,
		chsrcDir: chsrcDir,
	}
}

func (cm *ChsrcManager) GetName() string {
	return cm.Name
}

func (cm *ChsrcManager) IsInstalled() bool {
	// Check if chsrc is installed
	_, err := exec.LookPath("chsrc")
	return err == nil
}

func (cm *ChsrcManager) Installing(flags *GlobalFlags) error {
	// 获取最新的 GitHub 版本信息
	osType := "linux"
	if cm.config.OS == "darwin" {
		osType = "macos"
	}
	arch := "x64"
	if cm.config.ARCH == "arm64" {
		arch = "aarch64"
	}

	githubInfo := utils.NewGithubRepoInfo(
		"RubyMetric", "chsrc",
		fmt.Sprintf("chsrc-%s-%s", arch, osType),
		flags.HttpProxy,
		flags.GithubProxy,
	)

	var (
		tagName string
		err     error
	)
	if flags.Tag == "" {
		tagName, err = githubInfo.GetLatestReleaseTag()
		if err != nil {
			cm.config.Logger.Errorf("获取chsrc最新版本失败:%s", err)
			tagName = "v0.1.8"
		}
		cm.config.Logger.Infof("获取chsrc最新版本:%s", tagName)
	} else {
		tagName = flags.Tag
	}

	if err := os.MkdirAll(cm.chsrcDir, os.ModePerm); err != nil {
		cm.config.Logger.Errorf("创建目录%s失败:%s", cm.chsrcDir, err)
		return err
	}
	cm.config.Logger.Infof("已创建目录:%s", cm.chsrcDir)

	chsrcPath := fmt.Sprintf("%s/chsrc", cm.chsrcDir)
	if err := os.RemoveAll(chsrcPath); err != nil && !os.IsNotExist(err) {
		cm.config.Logger.Errorf("清理文件%s失败:%s", chsrcPath, err)
		return err
	}
	cm.config.Logger.Infof("清理文件%s", chsrcPath)

	downloadFile := fmt.Sprintf("%s/%s", cm.chsrcDir, "chsrc")
	cm.config.Logger.Infof("开始下载chsrc,版本:%s, 文件:%s", tagName, githubInfo.FileName)
	if err := githubInfo.DownloadReleaseLatestFile(downloadFile, tagName); err != nil {
		cm.config.Logger.Errorf("下载chsrc文件%s失败:%s", githubInfo.FileName, err)
		return err
	}
	cm.config.Logger.Infof("已下载chsrc文件:%s", githubInfo.FileName)

	// 设置文件为可执行权限
	err = os.Chmod(downloadFile, 0o755) // 0755 是 Unix 文件权限中的常见可执行权限
	if err != nil {
		cm.config.Logger.Errorf("设置文件%s为可执行权限失败:%s", downloadFile, err)
		os.Exit(1)
	}

	cm.config.Logger.Infof("设置文件%s为可执行权限成功。", downloadFile)

	// 更新 .bashrc 和 .zshrc
	shellFiles := []string{
		fmt.Sprintf("%s/.bashrc", cm.config.HomeDir),
		fmt.Sprintf("%s/.zshrc", cm.config.HomeDir),
	}
	contentToAdd := `
export PATH=$PATH:~/.local/bin
`
	for _, shellFile := range shellFiles {
		if err := utils.UpdateConfigFiles(shellFile, contentToAdd); err != nil {
			cm.config.Logger.Errorf("文件%s添加配置失败:%s", shellFile, err)
		} else {
			cm.config.Logger.Infof("文件%s添加配置成功", shellFile)
		}
	}
	return nil
}

func (cm *ChsrcManager) Install(flags *GlobalFlags) error {
	if !flags.Force && cm.IsInstalled() {
		cm.config.Logger.Warn("chsrc已经安装。使用 -f 选项强制重新安装。")
		return nil
	}

	// Add installation logic here
	cm.config.Logger.Info("开始安装chsrc...")
	if err := cm.Installing(flags); err != nil {
		os.Exit(1)
	}
	cm.config.Logger.Infof("chsrc安装成功!")
	return nil
}

func (cm *ChsrcManager) Update(flags *GlobalFlags) error {
	if !cm.IsInstalled() {
		cm.config.Logger.Warn("chsrc尚未安装。请使用 'install' 命令首先安装它。")
		return nil
	}

	// Add update logic here
	cm.config.Logger.Info("更新chsrc...")
	if err := cm.Installing(flags); err != nil {
		os.Exit(1)
	}
	cm.config.Logger.Infof("chsrc更新成功!")
	return nil
}

func (cm *ChsrcManager) Delete(flags *GlobalFlags) error {
	cm.config.Logger.Info("开始删除chsrc...")
	// 删除chsrc 二进制文件
	downloadFile := fmt.Sprintf("%s/%s", cm.chsrcDir, "chsrc")
	if err := os.Remove(downloadFile); err != nil && !os.IsNotExist(err) {
		cm.config.Logger.Errorf("删除chsrc可执行文件%s失败:%s", downloadFile, err)
		return err
	}
	cm.config.Logger.Infof("删除chsrc可执行文件%s成功", downloadFile)
	cm.config.Logger.Infof("chsrc删除成功!")
	return nil
}
