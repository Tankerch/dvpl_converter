package utils

import (
	"os"
	"path/filepath"
)

func IsDVPLFile(filename string) bool {
	return filepath.Ext(filename) == ".dvpl"
}

func IsProgramFile(filename string) bool {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}
	return absPath == exePath
}
