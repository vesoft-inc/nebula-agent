package utils

import (
	"os"
)

func IsFilePathExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
