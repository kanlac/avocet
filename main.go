package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "make",
		Usage: "",
		Action: func(cCtx *cli.Context) error {
			if isNotDirectory(cCtx.Args().Get(0)) {
				return fmt.Errorf("please input a valid directory path")
			}
			fmt.Println("good")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// isNotDirectory returns true if given path is not a valid directory path
func isNotDirectory(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return true
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return true
	}
	return !fileInfo.IsDir()
}
