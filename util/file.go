package util

import (
	"fmt"
	"os"
)

func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else {
		return false
	}
}

func CheckDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// create folder
			err := os.MkdirAll(dirPath, 0700)
			if err != nil {
				return fmt.Errorf("创建目录失败: %v", err)
			} else {
				return nil
			}
		} else {
			return err
		}
	} else {
		return nil
	}
}
