package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/bitfield/script"
	git "github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
)

type GithubRepoInfo struct {
	ower        string
	repo        string
	httpsProxy  string
	githubProxy string
	httpClinet  *http.Client
	logger      *logrus.Logger
}

func generateHttpClient(httpsProxy string) *http.Client {
	if httpsProxy == "" {
		return &http.Client{}
	}
	proxyURL, err := url.Parse(httpsProxy) // 设置代理地址
	if err != nil {
		return &http.Client{}
	}

	// 创建自定义的 Transport 以使用代理
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// 创建使用自定义 Transport 的 HTTP 客户端
	return &http.Client{
		Transport: transport,
	}
}

func NewGithubRepoInfo(ower, repo, httpsProxy, githubProxy string, logger *logrus.Logger) *GithubRepoInfo {
	return &GithubRepoInfo{
		ower:        ower,
		repo:        repo,
		httpsProxy:  httpsProxy,
		githubProxy: githubProxy,
		httpClinet:  generateHttpClient(httpsProxy),
		logger:      logger,
	}
}

func (g *GithubRepoInfo) GetLatestReleaseTag() string {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", g.ower, g.repo)
	tagName, err := script.NewPipe().WithHTTPClient(g.httpClinet).Get(api).JQ(".tag_name").String()
	if err != nil {
		g.logger.Errorf("获取%s最新版本失败:%s", g.repo, err)
		return ""
	}
	if len(tagName) < 2 {
		g.logger.Errorf("获取%s最新版本%s有误", g.repo, tagName)
		return ""
	}
	n := len(tagName)
	g.logger.Infof("获取%s最新版本:%s", g.repo, tagName)
	return tagName[1 : n-2]
}

func (g *GithubRepoInfo) DownloadReleaseLatestFile(dstFileName, srcFileName, tagName string) error {
	downloadUrl := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		g.ower, g.repo, tagName, srcFileName,
	)
	g.logger.Infof("开始从repo: %s 的Release下载 %s %s", g.repo, tagName, srcFileName)
	if g.githubProxy != "" {
		downloadUrl = JoinURL(g.githubProxy, downloadUrl)
	}
	_, err := script.NewPipe().WithHTTPClient(g.httpClinet).Get(downloadUrl).WriteFile(dstFileName)
	if err != nil {
		g.logger.Errorf("文件%s下载失败:%s", srcFileName, err)
		return err
	}
	g.logger.Infof("文件%s下载成功", srcFileName)
	return nil
}

func (g *GithubRepoInfo) GetRepoUrl() string {
	repoUrl := fmt.Sprintf(
		"https://github.com/%s/%s.git",
		g.ower, g.repo,
	)
	if g.githubProxy != "" {
		repoUrl = JoinURL(g.githubProxy, repoUrl)
	}
	return repoUrl
}

func (g *GithubRepoInfo) CloneRepo(dstPath string) error {
	g.logger.Infof("检测本地路径:%s是否存在repo:%s...", dstPath, g.repo)
	if DirectoryExists(dstPath) {
		if _, err := git.PlainOpen(dstPath); err == nil {
			g.logger.Infof("本地路径:%s是一个正常的Git仓库,不需要重新 Clone repo:%s", dstPath, g.repo)
			return nil
		}
		g.logger.Infof("本地路径:%s已存在,但不是正常的Git仓库,需要先删除再重新Clone", dstPath)
		if err := os.RemoveAll(dstPath); err != nil && !os.IsNotExist(err) {
			g.logger.Errorf("本地路径:%s清理失败:%s", dstPath, err)
			return err
		}
	}
	g.logger.Infof("本地不存在repo:%s,需要从远程 Clone 到本地:%s", g.repo, dstPath)
	if _, err := git.PlainClone(dstPath, false, &git.CloneOptions{
		Depth:    1,
		URL:      g.GetRepoUrl(),
		Progress: os.Stdout,
	}); err != nil {
		g.logger.Errorf("Clone repo:%s失败:%s", g.repo, err)
		return err
	} else {
		g.logger.Infof("Clone repo:%s成功", g.repo)
	}

	return nil
}

func (g *GithubRepoInfo) PullRepo(dstPath string) error {
	g.logger.Infof("检测本地路径:%s是否存在repo:%s...", dstPath, g.repo)
	if !DirectoryExists(dstPath) {
		g.logger.Infof("本地路径:%s不存在, 无法pull repo:%s", dstPath, g.repo)
		return fmt.Errorf("本地路径:%s不存在", dstPath)
	}

	// Open the existing repository
	repo, err := git.PlainOpen(dstPath)
	if err != nil {
		g.logger.Infof("本地路径:%s存在, 但不是一个Git仓库, 无法pull repo:%s", dstPath, g.repo)
		return err
	}
	g.logger.Infof("Pulling最新变更到本地仓库:%s...", dstPath)

	// Get the working tree for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		g.logger.Errorf("本地仓库%s无法获取工作树, 错误: %s", dstPath, err)
		return err
	}

	err = worktree.Reset(&git.ResetOptions{Mode: git.HardReset})
	if err != nil {
		g.logger.Errorf("执行 git reset --hard 错误: %v", err)
	}
	g.logger.Infof("成功执行 git reset --hard")

	// 执行 git clean -d --force 操作
	err = worktree.Clean(&git.CleanOptions{
		Dir: true,
	})
	if err != nil {
		g.logger.Errorf("执行 git clean -d --force 错误: %v", err)
	}
	g.logger.Infof("成功执行 git clean -d --force")

	// Pull the latest changes from the remote repository
	err = worktree.Pull(&git.PullOptions{
		RemoteName:        "origin",
		Progress:          os.Stdout,
		Depth:             1,
		Force:             true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			g.logger.Infof("本地仓库%s已经是最新的,无需拉取", dstPath)
			return nil
		}
		g.logger.Errorf("执行 git pull --rebase 错误: %s", err)
		return err
	}
	g.logger.Infof("成功执行 git pull --rebase")

	g.logger.Infof("Pull repo:%s成功", g.repo)
	return nil
}
