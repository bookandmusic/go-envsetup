package app

import (
	"fmt"
	"os"

	"github.com/bookandmusic/envsetup/config"
	"github.com/bookandmusic/envsetup/utils"
)

// Define ChsrcManager to handle chsrc operations
type ChsrcManager struct {
	Name    string
	ower    string
	repo    string
	tagName string
	config  *config.Config
}

func NewChsrcManager() *ChsrcManager {
	config := config.GetConfig()
	return &ChsrcManager{
		Name:    "chsrc",
		ower:    "RubyMetric",
		repo:    "chsrc",
		tagName: "v0.1.8",
		config:  config,
	}
}

func (cm *ChsrcManager) GetName() string {
	return cm.Name
}

func (cm *ChsrcManager) isInstalled() bool {
	return utils.IsCommandAvailable("chsrc")
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
	srcFileName := fmt.Sprintf("chsrc-%s-%s", arch, osType)

	githubInfo := utils.NewGithubRepoInfo(
		cm.ower, cm.repo,
		flags.HttpProxy,
		flags.GithubProxy,
		cm.config.Logger,
	)

	var tagName string
	if flags.Tag == "" {
		tagName = githubInfo.GetLatestReleaseTag()
		if tagName == "" {
			tagName = cm.tagName
		}
	} else {
		tagName = flags.Tag
	}

	downloadFile := fmt.Sprintf("/tmp/%s", "chsrc")
	if err := utils.RemoveFile(downloadFile, cm.config.Logger); err != nil {
		return err
	}

	if err := githubInfo.DownloadReleaseLatestFile(downloadFile, srcFileName, tagName); err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("install -m 755 %s /usr/local/bin", downloadFile)
	cmdStr = utils.GenerateCmd(cmdStr, true, cm.config.IsRoot)
	if err := utils.ExecCmd(cmdStr, cm.config.Logger); err != nil {
		return err
	}

	return nil
}

func (cm *ChsrcManager) Install(flags *GlobalFlags) error {
	if !flags.Force && cm.isInstalled() {
		cm.config.Logger.Warn("chsrc已经安装。使用 -f 选项强制重新安装。")
		return nil
	}

	// Add installation logic here
	cm.config.Logger.Info("开始安装chsrc...")
	if err := cm.Installing(flags); err != nil {
		cm.config.Logger.Errorf("chsrc安装失败!")
		os.Exit(1)
	}
	cm.config.Logger.Infof("chsrc安装成功!")
	return nil
}

func (cm *ChsrcManager) Update(flags *GlobalFlags) error {
	if !cm.isInstalled() {
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
	cmdStr := "rm -rf $(which chsrc)"
	cmdStr = utils.GenerateCmd(cmdStr, true, cm.config.IsRoot)
	if err := utils.ExecCmd(cmdStr, cm.config.Logger); err != nil {
		cm.config.Logger.Errorf("chsrc删除失败!")
		return err
	}
	cm.config.Logger.Infof("chsrc删除成功!")
	return nil
}
