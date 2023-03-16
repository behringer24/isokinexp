/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	srcPath  string
	destPath string

	rootCmd = &cobra.Command{
		Use:     "isokinexp",
		Short:   "Exporter / renamer for isokinetic device files",
		Version: "0.0.2",
		Long:    `isokinexp is a command to import export files from isokinetic measureing devices.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			copy()
		},
	}
)

func copy() {
	fmt.Printf("Copy files from %s to %s\n", srcPath, destPath)

	// Read all files from the source directory
	files, err := ioutil.ReadDir(srcPath)
	if err != nil {
		panic(err)
	}

	// Define regular expression patterns to match the filename and date information
	datePattern := regexp.MustCompile(`Date of Test:\s*(\d{2})\.(\d{2})\.(\d{4})`)
	namePattern := regexp.MustCompile(`Name of Person:\s*(\w+)`)

	// Iterate over each file in the source directory
	for _, file := range files {
		// Check if the file is a .txt file
		if !file.IsDir() && filepath.Ext(file.Name()) == ".585" {
			// Read the contents of the file
			content, err := ioutil.ReadFile(filepath.Join(srcPath, file.Name()))
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
				continue
			}

			// Extract the filename and date information from the contents of the file
			var name, date string
			for _, line := range strings.Split(string(content), "\n") {
				if matches := datePattern.FindStringSubmatch(line); len(matches) > 1 {
					date = matches[3] + "-" + matches[2] + "-" + matches[1]
				}
				if matches := namePattern.FindStringSubmatch(line); len(matches) > 1 {
					name = matches[1]
				}
				if name != "" && date != "" {
					break
				}
			}

			if name == "" || date == "" {
				fmt.Printf("File %s does not contain valid name and date information\n", file.Name())
				continue
			}

			// Create the destination directory based on the filename and date information
			destSubdir := filepath.Join(destPath, date)
			err = os.MkdirAll(destSubdir, 0755)
			if err != nil {
				fmt.Printf("Error creating destination directory %s: %v\n", destSubdir, err)
				continue
			}

			// Rename and move the file to the destination directory
			newName := fmt.Sprintf("%s_%s.txt", name, strings.ReplaceAll(date, "-", ""))
			err = os.Rename(filepath.Join(srcPath, file.Name()), filepath.Join(destSubdir, newName))
			if err != nil {
				fmt.Printf("Error moving file %s: %v\n", file.Name(), err)
				continue
			}

			fmt.Printf("File %s moved to %s\n", file.Name(), filepath.Join(destSubdir, newName))
		} else {
			fmt.Printf("Skip %s\n", file.Name())
		}
	}
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.isokinexp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&srcPath, "in", "i", "", "path to import files from")
	rootCmd.Flags().StringVarP(&destPath, "out", "o", "", "path to export files to")
	rootCmd.Flags().BoolP("delete", "d", false, `delete source files from input directory.
isokinexp will move files instead of copying`)

	rootCmd.MarkFlagRequired("in")
	rootCmd.MarkFlagRequired("out")
}
