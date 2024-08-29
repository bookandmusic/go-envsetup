package utils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/bitfield/script"
)

type GithubRepoInfo struct {
	Ower        string
	Repo        string
	FileName    string
	HttpsProxy  string
	GithubProxy string
}

var httpClinet *http.Client = nil

func setHttpClient(httpsProxy string) error {
	if httpClinet != nil {
		return nil
	}
	if httpsProxy == "" {
		httpClinet = &http.Client{}
		return nil
	}
	proxyURL, err := url.Parse(httpsProxy) // 设置代理地址
	if err != nil {
		httpClinet = &http.Client{}
		return nil
	}

	// 创建自定义的 Transport 以使用代理
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// 创建使用自定义 Transport 的 HTTP 客户端
	httpClinet = &http.Client{
		Transport: transport,
	}
	return nil
}

func NewGithubRepoInfo(ower, repo, filename, httpsProxy, githubProxy string) *GithubRepoInfo {
	setHttpClient(httpsProxy)
	return &GithubRepoInfo{
		Ower:        ower,
		Repo:        repo,
		FileName:    filename,
		HttpsProxy:  httpsProxy,
		GithubProxy: githubProxy,
	}
}

func (g *GithubRepoInfo) GetLatestReleaseTag() (string, error) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", g.Ower, g.Repo)
	tagName, err := script.NewPipe().WithHTTPClient(httpClinet).Get(api).JQ(".tag_name").String()
	if err != nil || len(tagName) < 2 {
		return "", fmt.Errorf("failed to get latest release: %v", err)
	}
	n := len(tagName)
	return tagName[1 : n-2], nil
}

func (g *GithubRepoInfo) DownloadReleaseLatestFile(path, tagName string) error {
	downloadUrl := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		g.Ower, g.Repo, tagName, g.FileName,
	)
	if g.GithubProxy != "" {
		downloadUrl = JoinURL(g.GithubProxy, downloadUrl)
	}
	_, err := script.NewPipe().WithHTTPClient(httpClinet).Get(downloadUrl).WriteFile(path)
	if err != nil {
		return err
	}
	return nil
}
