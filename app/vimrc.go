package app

import (
	"fmt"

	"github.com/bookandmusic/envsetup/config"
	"github.com/bookandmusic/envsetup/utils"
)

type VimrcManager struct {
	Name     string
	ower     string
	repo     string
	config   *config.Config
	vimrcDir string
}

func NewVimrcManager() *VimrcManager {
	config := config.GetConfig()

	vimrcDir := fmt.Sprintf("%s/.vim_runtime", config.HomeDir)
	return &VimrcManager{
		Name:     "vimrc",
		ower:     "amix",
		repo:     "vimrc",
		config:   config,
		vimrcDir: vimrcDir,
	}
}

func (v *VimrcManager) GetName() string {
	return v.Name
}

func (v *VimrcManager) Install(flags *GlobalFlags) error {
	installer, err := utils.GetInstaller(v.config.IsRoot, v.config.Logger)
	if err != nil {
		v.config.Logger.Errorf(err.Error())
		return err
	}
	if err := installer.CheckInstall("vim", "vim"); err != nil {
		return err
	}
	githubInfo := utils.NewGithubRepoInfo(
		v.ower, v.repo,
		flags.HttpProxy,
		flags.GithubProxy,
		v.config.Logger,
	)

	if err := githubInfo.CloneRepo(v.vimrcDir); err != nil {
		return err
	}
	cmdStr := "sh ~/.vim_runtime/install_awesome_vimrc.sh"
	if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
		v.config.Logger.Errorf("vimrv安装失败!")
		return err
	}
	v.config.Logger.Infof("vimrv安装成功!")
	return nil
}

func (v *VimrcManager) Update(flags *GlobalFlags) error {
	githubInfo := utils.NewGithubRepoInfo(
		v.ower, v.repo,
		flags.HttpProxy,
		flags.GithubProxy,
		v.config.Logger,
	)

	if err := githubInfo.PullRepo(v.vimrcDir); err != nil {
		return err
	}

	var cmdStr string
	updatePluginFile := "~/.vim_runtime/update_plugins.py"
	tmpUpdatePluginFile := "~/.vim_runtime/update_plugins-bak.py"
	if flags.GithubProxy != "" {
		cmdStr = fmt.Sprintf("cp %s %s && sed -i 's|https://github.com|https://mirror.ghproxy.com/https://github.com|g' %s", updatePluginFile, tmpUpdatePluginFile, updatePluginFile)
		if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
			v.config.Logger.Errorf("更新GitHub镜像地址失败!")
			return err
		}
		v.config.Logger.Infof("更新GitHub镜像地址成功!")
	}

	cmdStr = ""
	if flags.HttpProxy != "" {
		cmdStr = fmt.Sprintf("export https_proxy=%s && ", flags.HttpProxy)
	}
	if utils.IsCommandAvailable("python") {
		cmdStr = cmdStr + "python ~/.vim_runtime/update_plugins.py"
	} else if utils.IsCommandAvailable("python3") {
		cmdStr = cmdStr + "python3 ~/.vim_runtime/update_plugins.py"
	} else {
		v.config.Logger.Errorf("系统中不存在python解释器，无法更新插件")
		return nil
	}
	if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
		v.config.Logger.Errorf("更新插件失败!")
		return err
	}
	v.config.Logger.Infof("更新插件成功!")
	if flags.GithubProxy != "" {
		cmdStr = fmt.Sprintf("cp %s %s && rm -rf %s", tmpUpdatePluginFile, updatePluginFile, tmpUpdatePluginFile)
		if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
			v.config.Logger.Errorf("删除临时插件文件失败!")
			return err
		}
		v.config.Logger.Infof("删除临时插件文件成功!")
	}
	return nil
}

func (v *VimrcManager) Delete(flags *GlobalFlags) error {
	v.config.Logger.Info("开始删除vimrc...")

	cmdStr := "rm -rf ~/.vim_runtime ~/.vimrc"
	if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
		v.config.Logger.Errorf("删除~/.vim_runtime和~/.vimrc失败!")
		return err
	}
	v.config.Logger.Infof("删除~/.vim_runtime和~/.vimrc成功!")
	return nil
}
