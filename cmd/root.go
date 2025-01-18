package cmd

import (
	"fmt"
	"io"
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
	moveFile bool
	filter   string

	rootCmd = &cobra.Command{
		Use:     "isokinexp",
		Short:   "Exporter / renamer for isokinetic device files",
		Version: "0.2.1",
		Long:    `isokinexp is a command to import export files from isokinetic measureing devices.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			copy()
		},
	}
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func copy() {
	var count int = 0

	if moveFile {
		fmt.Printf("Move files from %s to %s\n", srcPath, destPath)
	} else {
		fmt.Printf("Copy files from %s to %s\n", srcPath, destPath)
	}

	// Read all files from the source directory
	files, err := ioutil.ReadDir(srcPath)
	if err != nil {
		panic(err)
	}

	// Define regular expression patterns to match the filename and date information
	datePattern := regexp.MustCompile(`Date of Test:\s*(\d{2})\.(\d{2})\.(\d{4})`)
	timePattern := regexp.MustCompile(`Time of Test:\s*(\d{2}\:\d{2})`)
	namePattern := regexp.MustCompile(`Name of Person:\s*(\w+)`)

	if filter != "" {
		fmt.Printf("Filter names by: %s\n", filter)
	}
	filterPattern := regexp.MustCompile(`^` + strings.Replace(filter, "*", ".*", -1) + `$`)
	filePattern := regexp.MustCompile(`\.\d{3}`)
	fmt.Printf("\n")

	// Iterate over each file in the source directory
	for _, file := range files {
		// Check if the file is not a directory and has a 3 digit feliname ext
		if !file.IsDir() && filePattern.Match([]byte(filepath.Ext(file.Name()))) {
			// Read the contents of the file
			content, err := ioutil.ReadFile(filepath.Join(srcPath, file.Name()))
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
				continue
			}

			// Extract the filename and date information from the contents of the file
			var name, date, time string
			for _, line := range strings.Split(string(content), "\n") {
				if matches := datePattern.FindStringSubmatch(line); len(matches) > 1 {
					date = matches[3] + "-" + matches[2] + "-" + matches[1]
				}
				if matches := namePattern.FindStringSubmatch(line); len(matches) > 1 {
					name = matches[1]
				}
				if matches := timePattern.FindStringSubmatch(line); len(matches) > 1 {
					time = strings.Replace(matches[1], ":", "-", 1)
				}
				if name != "" && date != "" {
					break
				}
			}

			if name == "" || date == "" || time == "" {
				fmt.Printf("File %s does not contain valid name, date and time information\n", file.Name())
				continue
			}

			if filter != "" && !filterPattern.Match([]byte(name)) {
				fmt.Printf("File %s skipped '%s' does not match '%s'\n", file.Name(), name, filter)
				continue
			}

			// Create the destination directory based on the filename and date information
			destSubdir := filepath.Join(destPath, name)
			err = os.MkdirAll(destSubdir, 0755)
			if err != nil {
				fmt.Printf("Error creating destination directory %s: %v\n", destSubdir, err)
				continue
			}

			// Rename and move the file to the destination directory
			newName := filepath.Join(destSubdir, fmt.Sprintf("%s_%s-%s.txt", name, date, time))

			if moveFile {
				err = os.Rename(filepath.Join(srcPath, file.Name()), newName)
				if err != nil {
					fmt.Printf("Error moving file %s: %v\n", file.Name(), err)
					continue
				}
				fmt.Printf("File %s moved to %s\n", file.Name(), newName)
				count++
			} else {
				err = copyFile(filepath.Join(srcPath, file.Name()), newName)
				if err != nil {
					fmt.Printf("Error copying file %s: %v\n", file.Name(), err)
					continue
				}
				fmt.Printf("File %s copied to %s\n", file.Name(), newName)
				count++
			}
		} else {
			fmt.Printf("Skip %s\n", file.Name())
		}
	}
	fmt.Printf("%d files copied/moved\n", count)
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
	rootCmd.Flags().StringVarP(&filter, "filter", "f", "", "filter names (use * as placeholder)")
	rootCmd.Flags().BoolVarP(&moveFile, "delete", "d", false, `delete source files from input directory.
isokinexp will move files instead of copying`)

	rootCmd.MarkFlagRequired("in")
	rootCmd.MarkFlagRequired("out")
}
