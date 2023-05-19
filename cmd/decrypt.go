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

func convertDVPLtoFile(path string) {
	// Validation
	fileIsDVPL := utils.IsDVPLFile(path)
	if !fileIsDVPL {
		return
	}

	// Input
	fileBuf, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Processed
	outputBuf, err := dvpl.DecryptDVPL(fileBuf)
	if err != nil {
		panic(err)
	}

	// Output
	var outputPath = utils.DVPLOriginalName(path)
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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decompress all .dvpl to respective format",
	Long: `Decompress all .dvpl to respective format

Example:
- dvpl_converter decrypt (Decompress all files, included subfolder in current directory)
- dvpl_converter decrypt -d 3d (Decompress all files, included subfolder in 3d's folder)
- dvpl_converter decrypt -d 3d --delete-original (Delete original file after decompres all files in 3d's folder)
- dvpl_converter decrypt -d 3d/Tanks/France/Images/B1.mali.pvr.dvpl (Decompress only this file)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start decrypting:")
		dirInfo, err := os.Stat(inputDirPath)
		if err != nil {
			panic(err)
		}

		// dirPath is single file
		if !dirInfo.IsDir() {
			convertDVPLtoFile(inputDirPath)
			return
		}

		// dirPath is Directory, included default value (cwd)
		filepath.WalkDir(inputDirPath, func(path string, d fs.DirEntry, ___ error) error {
			if d.IsDir() {
				return nil
			}
			convertDVPLtoFile(path)
			return nil
		})

	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVarP(&inputDirPath, "dir", "d", ".", "Specify the input directory/file path, The path can be absolute or relative")
	decryptCmd.Flags().BoolVar(&deleteOriginalFlag, "delete-original", false, `Delete original file after conversion. Warning: This action is irreversible.`)
}
