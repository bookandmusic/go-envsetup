package utils

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Installer interface defines the required methods for all installers.
type Installer interface {
	GetIsRoot() bool
	InputSudoPasswd() error
	GetPackageManager() string
	Install(packages []string) error
	Unintstall(packages []string) error
	CheckUnInstall(name, command string) error
	GetIsRequiresSudo() bool
	CheckInstall(name, command string) error
}

// BaseInstaller struct, holds common methods and properties for all installers.
type BaseInstaller struct {
	packageManager string
	isRequiresSudo bool
	isRoot         bool
	logger         *logrus.Logger
}

func (bi *BaseInstaller) GetPackageManager() string {
	return bi.packageManager
}

func (bi *BaseInstaller) GetIsRoot() bool {
	return bi.isRoot
}

func (bi *BaseInstaller) GetIsRequiresSudo() bool {
	return bi.isRequiresSudo
}

func (bi *BaseInstaller) InputSudoPasswd() error {
	cmd := exec.Command("sudo", "-v")
	return cmd.Run()
}

func (bi *BaseInstaller) Install(packages []string) error {
	cmdStr := fmt.Sprintf("%s install -y", bi.packageManager)
	for _, pkg := range packages {
		cmdStr += fmt.Sprintf(" %s", pkg)
	}
	cmdStr = GenerateCmd(cmdStr, true, bi.isRoot)

	return ExecCmd(cmdStr, bi.logger)
}

func (bi *BaseInstaller) Unintstall(packages []string) error {
	cmdStr := fmt.Sprintf("%s uninstall -y", bi.packageManager)
	for _, pkg := range packages {
		cmdStr += fmt.Sprintf(" %s", pkg)
	}

	cmdStr = GenerateCmd(cmdStr, true, bi.isRoot)

	return ExecCmd(cmdStr, bi.logger)
}

func (bi *BaseInstaller) CheckInstall(name, command string) error {
	if !IsCommandAvailable(command) {
		bi.logger.Warnf("检测%s不存在,需要先安装%s", name, name)
		if bi.isRequiresSudo && !bi.isRoot {
			bi.logger.Infof("当前用户不是root用户,且需要sudo权限,请输入sudo密码")
			bi.InputSudoPasswd()
		}
		bi.logger.Infof("开始使用%s安装%s...", bi.packageManager, name)
		if err := bi.Install([]string{name}); err != nil {
			bi.logger.Errorf("%s安装失败,%s", name, err)
			return err
		} else {
			bi.logger.Infof("%s安装成功!!!", name)
		}
	} else {
		bi.logger.Infof("检测%s已存在", name)
	}
	return nil
}

func (bi *BaseInstaller) CheckUnInstall(name, command string) error {
	if !IsCommandAvailable(command) {
		bi.logger.Warnf("检测%s不存在,不需要卸载", name)
		return nil
	}
	bi.logger.Infof("检测%s存在", name)
	if bi.isRequiresSudo && !bi.isRoot {
		bi.logger.Infof("当前用户不是root用户,且需要sudo权限,请输入sudo密码")
		bi.InputSudoPasswd()
	}
	bi.logger.Infof("开始使用%s卸载%s...", bi.packageManager, name)
	if err := bi.Unintstall([]string{name}); err != nil {
		bi.logger.Errorf("%s卸载失败,%s", name, err)
		return err
	} else {
		bi.logger.Infof("%s卸载成功!!!", name)
	}
	return nil
}

// AptInstaller struct for handling apt-get specific installation.
type AptInstaller struct {
	BaseInstaller
}

func (apt *AptInstaller) Unintstall(packages []string) error {
	cmdStr := fmt.Sprintf("%s remove -y --purge", apt.packageManager)
	for _, pkg := range packages {
		cmdStr += fmt.Sprintf(" %s", pkg)
	}
	apt.logger.Infof(cmdStr)

	cmdStr = GenerateCmd(cmdStr, true, apt.isRoot)

	return ExecCmd(cmdStr, apt.logger)
}

// NewAptInstaller creates a new AptInstaller instance.
func NewAptInstaller(isRoot bool, logger *logrus.Logger) *AptInstaller {
	return &AptInstaller{
		BaseInstaller{
			packageManager: "apt-get",
			isRequiresSudo: true,
			isRoot:         isRoot,
			logger:         logger,
		},
	}
}

// YumInstaller struct for handling yum specific installation.
type YumInstaller struct {
	BaseInstaller
}

// NewYumInstaller creates a new YumInstaller instance.
func NewYumInstaller(isRoot bool, logger *logrus.Logger) *YumInstaller {
	return &YumInstaller{
		BaseInstaller{
			packageManager: "yum",
			isRequiresSudo: true,
			isRoot:         isRoot,
			logger:         logger,
		},
	}
}

// BrewInstaller struct for handling Homebrew specific installation.
type BrewInstaller struct {
	BaseInstaller
}

// NewBrewInstaller creates a new BrewInstaller instance.
func NewBrewInstaller(isRoot bool, logger *logrus.Logger) *BrewInstaller {
	return &BrewInstaller{
		BaseInstaller{
			packageManager: "brew",
			isRequiresSudo: false,
			isRoot:         isRoot,
			logger:         logger,
		},
	}
}

// PortInstaller struct for handling MacPorts specific installation.
type PortInstaller struct {
	BaseInstaller
}

// NewPortInstaller creates a new PortInstaller instance.
func NewPortInstaller(isRoot bool, logger *logrus.Logger) *PortInstaller {
	return &PortInstaller{
		BaseInstaller{
			packageManager: "port",
			isRequiresSudo: true,
			isRoot:         isRoot,
			logger:         logger,
		},
	}
}

// GetInstaller returns an appropriate installer instance based on available package manager.
func GetInstaller(isRoot bool, logger *logrus.Logger) (Installer, error) {
	if IsCommandAvailable("apt-get") {
		return NewAptInstaller(isRoot, logger), nil
	} else if IsCommandAvailable("yum") {
		return NewYumInstaller(isRoot, logger), nil
	} else if IsCommandAvailable("brew") {
		// Homebrew does not require sudo
		return NewBrewInstaller(isRoot, logger), nil
	} else if IsCommandAvailable("port") {
		return NewPortInstaller(isRoot, logger), nil
	} else {
		return nil, fmt.Errorf("找不到适合的包管理器")
	}
}
