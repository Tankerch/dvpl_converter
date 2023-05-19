package utils

import (
	"strings"
)

func DVPLOriginalName(filename string) string {
	splitName := strings.Split(filename, ".")
	return strings.Join(splitName[:len(splitName)-1], ".")
}
