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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	decryptCmd.Flags().StringVarP(&inputDirPath, "dir", "d", ".", "Help message for toggle")
	decryptCmd.Flags().BoolVar(&deleteOriginalFlag, "delete-original", false, "Help message for toggle")
}
