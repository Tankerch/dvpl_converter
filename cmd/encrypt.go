/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tankerch/dvpl_converter/common/dvpl"
	"github.com/tankerch/dvpl_converter/common/utils"
)

func convertFileToDVPL(path string) {
	// Validation
	fileIsDVPL := utils.IsDVPLFile(path)
	if fileIsDVPL {
		return
	}

	// Input
	fileBuf, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Processed
	outputBuf, err := dvpl.EncryptDVPL(fileBuf)
	if err != nil {
		panic(err)
	}

	// Output
	var outputPath = fmt.Sprintf("%s.dvpl", path)
	fout, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	fout.Write(outputBuf)

	// (Optional) Delete original
	if deleteOriginalFlag {
		if err := os.Remove(path); err != nil {
			panic(err)
		}
	}
	fmt.Printf("\t%s\n", path)
}

func StartEcryption() {
	fmt.Println("Start Encrypting:")
	dirInfo, err := os.Stat(inputDirPath)
	if err != nil {
		panic(err)
	}

	// dirPath is single file
	if !dirInfo.IsDir() {
		convertFileToDVPL(inputDirPath)
		return
	}

	// dirPath is Directory, included default value (cwd)
	filepath.WalkDir(inputDirPath, func(path string, d fs.DirEntry, ___ error) error {
		if d.IsDir() {
			return nil
		}
		convertFileToDVPL(path)
		return nil
	})

}

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Compress all files to .dvpl",
	Long: `Compress all files to .dvpl

Example:
- dvpl_converter encrypt (Compress all files, included subfolder in current directory)
- dvpl_converter encrypt -d 3d (Compress all files, included subfolder in 3d's folder)
- dvpl_converter encrypt -d 3d --delete-original (Delete original file after compres all files in 3d's folder)
- dvpl_converter encrypt -d 3d/Tanks/France/Images/B1.mali.pvr.dvpl (Compress only this file)`,
	Run: func(cmd *cobra.Command, args []string) {
		StartEcryption()
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVarP(&inputDirPath, "dir", "d", ".", "Specify the input directory/file path, The path can be absolute or relative")
	encryptCmd.Flags().BoolVar(&deleteOriginalFlag, "delete-original", false, `Delete original file after conversion. Warning: This action is irreversible.`)
}
