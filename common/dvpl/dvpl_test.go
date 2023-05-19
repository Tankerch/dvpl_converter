package dvpl

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/exp/slices"
)

func TestEncryptdecrypt(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(b), "../..")

	testFiles := []string{"./test_files/M3_Lee.mali.pvr", "./test_files/ARL_44.tex"}

	for _, path := range testFiles {
		originalBuf, err := os.ReadFile(filepath.Join(rootDir, path))
		if err != nil {
			panic(err)
		}
		compressedBuf, err := EncryptDVPL(originalBuf)
		if err != nil {
			panic(err)
		}

		decompressedBuf, err := DecryptDVPL(compressedBuf)
		if err != nil {
			panic(err)
		}

		if !slices.Equal(originalBuf, decompressedBuf) {
			t.Errorf("Original and decompressed result should same")
		}

	}
}
