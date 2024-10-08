package app

import (
	"fmt"
	"os"

	archiver "github.com/mholt/archiver/v3"

	"github.com/bookandmusic/envsetup/config"
	"github.com/bookandmusic/envsetup/utils"
)

// Define VMRManager to handleVMRoperations
type VMRManager struct {
	Name    string
	ower    string
	repo    string
	tagName string
	config  *config.Config
	vmrDir  string
}

func NewVMRManager() *VMRManager {
	config := config.GetConfig()

	vmrDir := fmt.Sprintf("%s/.vmr", config.HomeDir)
	return &VMRManager{
		Name:    "vmr",
		ower:    "gvcgo",
		tagName: "v0.6.5",
		repo:    "version-manager",
		config:  config,
		vmrDir:  vmrDir,
	}
}

func (vm *VMRManager) Installing(flags *GlobalFlags) error {
	srcFileName := fmt.Sprintf("vmr_%s-%s.zip", vm.config.OS, vm.config.ARCH)

	// 获取最新的 GitHub 版本信息
	githubInfo := utils.NewGithubRepoInfo(
		vm.ower, vm.repo,
		flags.HttpProxy,
		flags.GithubProxy,
		vm.config.Logger,
	)
	var tagName string
	if flags.Tag == "" {
		tagName = githubInfo.GetLatestReleaseTag()
		if tagName == "" {
			tagName = vm.tagName
		}
	} else {
		tagName = flags.Tag
	}

	if err := utils.Mkdir(vm.vmrDir, vm.config.Logger); err != nil {
		return err
	}

	downloadFile := fmt.Sprintf("%s/%s", vm.vmrDir, srcFileName)

	if err := utils.RemoveFile(downloadFile, vm.config.Logger); err != nil {
		return err
	}

	if err := githubInfo.DownloadReleaseLatestFile(downloadFile, srcFileName, tagName); err != nil {
		return err
	}

	// 使用 archiver 解压 ZIP 文件
	vmrPath := fmt.Sprintf("%s/vmr", vm.vmrDir)
	if err := utils.RemoveFile(vmrPath, vm.config.Logger); err != nil {
		return err
	}
	if err := archiver.Unarchive(downloadFile, vm.vmrDir); err != nil {
		vm.config.Logger.Errorf("解压VMR文件%s失败:%s", downloadFile, err)
		return err
	}
	vm.config.Logger.Infof("已解压VMR文件:%s", downloadFile)

	// 删除下载的VMR压缩文件
	if err := utils.RemoveFile(downloadFile, vm.config.Logger); err != nil {
		return nil
	}

	confPath := fmt.Sprintf("%s/conf.toml", vm.vmrDir)
	vm.config.Logger.Infof("生成VMR配置:%s", confPath)
	vmrConf := fmt.Sprintf(`
ProxyUri = ''
ReverseProxy = ''
SDKIntallationDir = '%s'
VersionHostUrl = 'https://gitee.com/moqsien/vsources/raw/main'
ThreadNum = 0
UseCustomedMirrors = true
`, vm.vmrDir)

	if err := os.WriteFile(confPath, []byte(vmrConf), 0o644); err != nil {
		vm.config.Logger.Errorf("生成VMR配置文件%s失败:%s", confPath, err)
		return err
	}

	mirrorsPath := fmt.Sprintf("%s/customed_mirrors.toml", vm.vmrDir)
	vm.config.Logger.Infof("生成VMR镜像配置:%s", mirrorsPath)
	customedMirrors := `
'https://go.dev/dl/' = 'https://mirrors.ustc.edu.cn/golang/'
'https://nodejs.org/download/release/' = 'https://mirrors.ustc.edu.cn/node/'
'https://repo.anaconda.com/miniconda/' = 'https://mirrors.ustc.edu.cn/anaconda/miniconda/'
`
	if err := os.WriteFile(mirrorsPath, []byte(customedMirrors), 0o644); err != nil {
		vm.config.Logger.Errorf("生成VMR镜像配置文件%s失败:%s", mirrorsPath, err)
		return err
	}

	scriptPath := fmt.Sprintf("%s/vmr.sh", vm.vmrDir)
	vm.config.Logger.Infof("生成VMR启动脚本:%s", scriptPath)
	vmrScript := fmt.Sprintf(`
# cd hook start
export PATH=%s:$PATH

if [ -z "$(alias|grep cdhook)" ]; then
	cdhook() {
		if [ $# -eq 0 ]; then
			cd
		else
			cd "$@" && vmr use -E
		fi
	}
	alias cd='cdhook'
fi

if [ ! $VMR_CD_INIT ]; then
        VMR_CD_INIT="vmr_cd_init"
        cd "$(pwd)"
fi
# cd hook end
`, vm.vmrDir)
	if err := os.WriteFile(scriptPath, []byte(vmrScript), 0o755); err != nil {
		vm.config.Logger.Errorf("生成VMR启动脚本%s失败:%s", scriptPath, err)
		return err
	}

	// 更新 .bashrc 和 .zshrc
	shellFiles := []string{
		fmt.Sprintf("%s/.bashrc", vm.config.HomeDir),
		fmt.Sprintf("%s/.zshrc", vm.config.HomeDir),
	}
	contentToAdd := `
# vm_envs start
if [ -z "$VM_DISABLE" ]; then
    . ~/.vmr/vmr.sh
fi
# vm_envs end
`
	for _, shellFile := range shellFiles {
		if !utils.FileExists(shellFile) {
			continue
		}
		if err := utils.UpdateConfigFiles(shellFile, contentToAdd); err != nil {
			vm.config.Logger.Errorf("文件%s添加配置失败：%s", shellFile, err)
		} else {
			vm.config.Logger.Infof("配置文件%s更新如下配置:", shellFile)
			vm.config.Logger.Infof(contentToAdd)
		}
	}
	return nil
}

func (vm *VMRManager) GetName() string {
	return vm.Name
}

func (vm *VMRManager) Install(flags *GlobalFlags) error {
	if !flags.Force && vm.isInstalled() {
		vm.config.Logger.Warn("VMR已经安装。使用 -f 选项强制重新安装。")
		return nil
	}

	vm.config.Logger.Infof("开始安装VMR...")
	if err := vm.Installing(flags); err != nil {
		os.Exit(1)
	}
	vm.config.Logger.Infof("VMR安装成功!")
	return nil
}

func (vm *VMRManager) Update(flags *GlobalFlags) error {
	if !vm.isInstalled() {
		vm.config.Logger.Warn("VMR尚未安装。请使用 'install' 命令首先安装它。")
		return nil
	}

	// Add update logic here
	vm.config.Logger.Info("更新VMR...")
	if err := vm.Installing(flags); err != nil {
		os.Exit(1)
	}
	vm.config.Logger.Infof("VMR更新成功!")
	return nil
}

func (vm *VMRManager) Delete(flags *GlobalFlags) error {
	// Add deletion logic here
	vm.config.Logger.Info("开始删除VMR...")

	// 删除VMR目录及其内容
	if err := os.RemoveAll(vm.vmrDir); err != nil {
		vm.config.Logger.Errorf("删除VMR目录%s失败:%s", vm.vmrDir, err)
		os.Exit(1)
	}
	vm.config.Logger.Infof("已删除VMR目录:%s", vm.vmrDir)

	// 从 .bashrc 和 .zshrc 中移除配置
	shellFiles := []string{
		fmt.Sprintf("%s/.bashrc", vm.config.HomeDir),
		fmt.Sprintf("%s/.zshrc", vm.config.HomeDir),
	}
	contentPattern := `# vm_envs start\nif \[ -z "\$VM_DISABLE" \]; then\n    \. ~/.vmr/vmr.sh\nfi\n# vm_envs end\n`
	for _, shellFile := range shellFiles {
		if err := utils.RemoveConfigFromFile(shellFile, contentPattern); err != nil {
			vm.config.Logger.Errorf("从文件%s移除配置失败:%s", shellFile, err)
			os.Exit(1)
		} else {
			vm.config.Logger.Infof("从文件%s移除配置成功", shellFile)
		}
	}

	vm.config.Logger.Infof("已从配置文件中移除VMR配置")
	vm.config.Logger.Infof("VMR删除成功!")
	return nil
}

func (vm *VMRManager) isInstalled() bool {
	return utils.IsCommandAvailable("vmr")
}
