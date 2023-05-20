package utils

import (
	"os"
	"path/filepath"
)

func IsDVPL(filename string) bool {
	return filepath.Ext(filename) == ".dvpl"
}

func IsProgramFile(filename string) bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return false
	}
	return absPath == exePath
}
