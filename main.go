package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "make",
		Usage: "",
		Action: func(cCtx *cli.Context) error {
			absPath, err := getAbsolutePath(cCtx.Args().Get(0))
			if err != nil {
				return fmt.Errorf("cannot get absolute path from %s: %+v", cCtx.Args().Get(0), err)
			}
			compressFiles(filepath.Base(absPath), getFilesInDirectory(absPath))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getAbsolutePath(pathArg string) (string, error) {
	if isNotDirectory(pathArg) {
		return "", fmt.Errorf("%s is not a directory", pathArg)
	}
	return filepath.Abs(pathArg)
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

// compressFiles compresses files and create a target.zip
func compressFiles(target string, files []string) {
	targetZip, err := openTargetZipFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed to open zip for writing: %+v", err)
	}
	defer targetZip.Close()

	zipWriter := zip.NewWriter(targetZip)
	defer zipWriter.Close()

	writeZipFiles(files, zipWriter)
	fmt.Println("ok")
}

// getFilesInDirectory returns file paths under dirname
// TODO
func getFilesInDirectory(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(files))
	return ret, nil
}

func openTargetZipFile(fileName string, flags int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(fileName+".zip", flags, perm)
}

func writeZipFiles(ssrcFilePaths []string, zipWriter *zip.Writer) {
	for _, filename := range ssrcFilePaths {
		if err := appendFileToZipWriter(filename, zipWriter); err != nil {
			log.Fatalf("Failed to add file %s to zip: %s", filename, err)
		}
	}
}

func appendFileToZipWriter(filename string, zipw *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {
		msg := "Failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Failed to write %s to zip: %s", filename, err)
	}

	return nil
}
