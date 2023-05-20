/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// Global flag
var inputDirPath string
var deleteOriginalFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dvpl_converter",
	Short: "Convert files to and from dvpl",
	Long: `dvpl_converter is a handy and reliable tool that lets you convert any file to dvpl format and vice versa.
	
.dvpl is a new file format that is first seen used in the World of Tanks Blitz Client for Chinese Server,
and now it's used on all known clients, except files that are contained within Android apks.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("failed to get current directory")
			return
		}

		qs := []*survey.Question{
			{
				Name: "action",
				Prompt: &survey.Select{
					Message: "Choose the disered action:",
					Options: []string{"encrypt", "decrypt"},
				},
			},
			{
				Name: "directory",
				Prompt: &survey.Input{
					Message: "Please enter the path of the directory/file you want to encrypt or decrypt: ",
					Default: cwd,
					Suggest: func(toComplete string) []string {
						files, _ := filepath.Glob(toComplete + "*")
						return files
					},
				},
			},
			{
				Name: "keepOriginal",
				Prompt: &survey.Confirm{
					Message: "Do you want to keep the original file after conversion?",
					Default: true,
				},
			},
		}
		answers := struct {
			Action       string
			Directory    string
			KeepOriginal bool
		}{}

		err = survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println("failed create survey prompt")
		}
		inputDirPath = answers.Directory
		deleteOriginalFlag = !answers.KeepOriginal

		// Execute
		if answers.Action == "encrypt" {
			StartEcryption()
		} else {
			StartDecrypting()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dvpl_converter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
