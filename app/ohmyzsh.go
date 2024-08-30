package app

import (
	"fmt"
	"time"

	"github.com/bookandmusic/envsetup/config"
	"github.com/bookandmusic/envsetup/utils"
)

type repo struct {
	ower      string
	repo      string
	localPath string
}

type OhMyZshManager struct {
	Name       string
	config     *config.Config
	ohMyZshDir string
	pluginDir  string
	themeDir   string
	repos      []*repo
}

func NewOhMyZshManager() *OhMyZshManager {
	config := config.GetConfig()
	ohMyZshDir := fmt.Sprintf("%s/.oh-my-zsh", config.HomeDir)
	pluginDir := fmt.Sprintf("%s/custom/plugins", ohMyZshDir)
	themeDir := fmt.Sprintf("%s/custom/themes", ohMyZshDir)
	return &OhMyZshManager{
		Name:       "ohmyzsh",
		config:     config,
		ohMyZshDir: ohMyZshDir,
		pluginDir:  pluginDir,
		themeDir:   themeDir,
		repos: []*repo{
			{
				ower:      "ohmyzsh",
				repo:      "ohmyzsh",
				localPath: ohMyZshDir,
			},
			{
				ower:      "zsh-users",
				repo:      "zsh-autosuggestions",
				localPath: fmt.Sprintf("%s/zsh-autosuggestions", pluginDir),
			},
			{
				ower:      "zsh-users",
				repo:      "zsh-syntax-highlighting",
				localPath: fmt.Sprintf("%s/zsh-syntax-highlighting", pluginDir),
			},
		},
	}
}

func (v *OhMyZshManager) GetName() string {
	return v.Name
}

func (v *OhMyZshManager) Install(flags *GlobalFlags) error {
	installer, err := utils.GetInstaller(v.config.IsRoot, v.config.Logger)
	if err != nil {
		v.config.Logger.Errorf(err.Error())
		return err
	}
	if err := installer.CheckInstall("zsh", "zsh"); err != nil {
		return err
	}

	for _, repo := range v.repos {
		githubInfo := utils.NewGithubRepoInfo(
			repo.ower, repo.repo,
			flags.HttpProxy,
			flags.GithubProxy,
			v.config.Logger,
		)

		if err := githubInfo.CloneRepo(repo.localPath); err != nil {
			return err
		}
	}

	zshrcPath := fmt.Sprintf("%s/.zshrc", v.config.HomeDir)
	if utils.FileExists(zshrcPath) {
		// 获取当前时间
		currentTime := time.Now()
		// 格式化时间为字符串 "yyyyMMddHHmmss"
		uniqueCode := currentTime.Format("20060102150405")
		v.config.Logger.Infof("本地%s已存在配置文件:.zshrc,将其备份为:.zshrc-%s", v.config.HomeDir, uniqueCode)
		cmdStr := fmt.Sprintf("cp ~/.zshrc ~/.zshrc-%s", uniqueCode)
		cmdStr = utils.GenerateCmd(cmdStr, false, v.config.IsRoot)
		if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
			v.config.Logger.Errorf("备份配置文件.zshrc失败!")
			return err
		}
		v.config.Logger.Infof("备份配置文件.zshrc成功!")
	}

	v.config.Logger.Infof("生成默认的配置文件: ~/.zshrc")
	cmdStr := "cp ~/.oh-my-zsh/templates/zshrc.zsh-template ~/.zshrc"
	if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
		v.config.Logger.Errorf("生成配置文件.zshrc失败!")
		return err
	}
	v.config.Logger.Infof("生成配置文件.zshrc成功!")

	cmdStr = "sed -i 's/plugins=(git)/plugins=(git sudo zsh-autosuggestions zsh-syntax-highlighting)/' ~/.zshrc"
	if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
		v.config.Logger.Errorf("修改~/.zshrc配置文件启用插件失败!")
		return err
	}
	v.config.Logger.Infof("修改~/.zshrc配置文件启用插件成功!")

	v.config.Logger.Infof("成功安装zsh及oh-my-zsh!!!")
	return nil
}

func (v *OhMyZshManager) Update(flags *GlobalFlags) error {
	if !utils.DirectoryExists(v.ohMyZshDir) {
		v.config.Logger.Warn("oh-my-zsh尚未安装。请使用 'install' 命令首先安装它。")
		return nil
	}
	for _, repo := range v.repos {
		githubInfo := utils.NewGithubRepoInfo(
			repo.ower, repo.repo,
			flags.HttpProxy,
			flags.GithubProxy,
			v.config.Logger,
		)

		if err := githubInfo.PullRepo(repo.localPath); err != nil {
			return err
		}
	}
	return nil
}

func (v *OhMyZshManager) Delete(flags *GlobalFlags) error {
	v.config.Logger.Info("开始删除ohmyzsh...")

	cmdStr := "rm -rf ~/.oh-my-zsh ~/.zshrc"
	if err := utils.ExecCmd(cmdStr, v.config.Logger); err != nil {
		v.config.Logger.Errorf("删除~/.oh-my-zsh和~/.zshrc失败!")
		return err
	}
	v.config.Logger.Infof("删除~/.oh-my-zsh和~/.zshrc成功!")
	return nil
}
