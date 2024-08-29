package utils

import (
	"bytes"
	"os"
	"regexp"
)

func UpdateConfigFiles(filePath, content string) error {
	fileContent, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !bytes.Contains(fileContent, []byte(content)) {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.WriteString(content); err != nil {
			return err
		}
	}
	return nil
}

// 从指定文件中移除配置内容
func RemoveConfigFromFile(filePath, contentPattern string) error {
	// 读取文件内容
	fileContent, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 创建正则表达式
	re, err := regexp.Compile(contentPattern)
	if err != nil {
		return err
	}

	// 替换内容
	newContent := re.ReplaceAllString(string(fileContent), "")

	// 写入新内容
	if err := os.WriteFile(filePath, []byte(newContent), 0o644); err != nil {
		return err
	}

	return nil
}
