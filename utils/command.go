package utils

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func GenerateCmd(cmdStr string, isSudo, isRoot bool) string {
	if !isSudo || isRoot {
		return cmdStr
	}
	return "sudo " + cmdStr
}

func ExecCmd(cmdStr string, logger *logrus.Logger) error {
	// 使用 exec.Command 创建命令
	cmd := exec.Command("bash", "-c", cmdStr)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行命令
	logger.Infof(cmdStr)
	return cmd.Run()
}
