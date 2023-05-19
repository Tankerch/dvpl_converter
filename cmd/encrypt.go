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

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		absName, err := filepath.Abs(inputDirPath)
		if err != nil {
			panic(err)
		}
		fmt.Printf("absname: %s | Dirpath: %s\n", absName, inputDirPath)
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

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVarP(&inputDirPath, "dir", "d", ".", "Help message for toggle")
	encryptCmd.Flags().BoolVar(&deleteOriginalFlag, "delete-original", false, "Help message for toggle")
}
