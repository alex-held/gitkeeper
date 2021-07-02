package main

import (
	"fmt"
	fs2 "io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	recursiveF = pflag.BoolP("recursive", "r", true, "gitkeeper -r .")
	recursive  bool
)

func main() {

	var addCmd = &cobra.Command{
		Use:   "add [OPTIONS] [path to add .gitkeep]",
		Short: "adds .gitkeep files to matching paths",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var absPath string
			if absPath, err = filepath.Abs(args[0]); err != nil {
				return err
			}

			println(absPath)

			var directories []string
			var emptyDirs []string

			fs := afero.NewOsFs()

			if !recursive {
				if dir, err := afero.IsDir(fs, absPath); err == nil && dir {
					directories = append(directories, absPath)
					if empty, err := afero.IsEmpty(fs, absPath); err == nil && empty {
						emptyDirs = append(emptyDirs, absPath)
					}
				}
				return err
			}

			err = afero.Walk(fs, absPath, func(path string, fi fs2.FileInfo, e error) (er error) {
				if e != nil {
					return e
				}
				if fi.IsDir() {
					directories = append(directories, path)
					if empty, er := afero.IsEmpty(fs, path); er == nil && empty {
						emptyDirs = append(emptyDirs, path)
					}
				}
				return er
			})

			fmt.Printf("Directories found: %d\n\n", len(directories))
			fmt.Printf("Empty Directories found: %d\n\n", len(emptyDirs))

			return err
		},
	}
	addCmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "gitkeeper -r .")

	var rootCmd = &cobra.Command{Use: "gitkeeper"}
	rootCmd.AddCommand(addCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
