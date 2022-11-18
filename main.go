package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.App{
		Name:  "make",
		Usage: "",
		Action: func(cCtx *cli.Context) error {
			arg := cCtx.Args().Get(0)
			if len(arg) == 0 {
				return fmt.Errorf("please provide a path argument")
			}
			absPath, err := getAbsolutePath(arg)
			if err != nil {
				return fmt.Errorf("cannot get absolute path from %s: %+v", arg, err)
			}
			var files []string
			if err = getFilesInDirectory(absPath, &files); err != nil {
				return fmt.Errorf("cannot get files in directory %s: %+v", absPath, err)
			}
			compressFiles(filepath.Base(absPath), files)
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

	writeZipFiles(files, zipWriter, target)
	fmt.Println("ok")
}

// getFilesInDirectory searches contents under dirname recursively
// and store the file paths in files.
func getFilesInDirectory(dirPath string, files *[]string) error {
	entires, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, e := range entires {
		if e.IsDir() {
			if err = getFilesInDirectory(filepath.Join(dirPath, e.Name()), files); err != nil {
				return err
			}
			continue
		}

		fi, err := e.Info()
		if err != nil {
			return err
		}
		*files = append(*files, filepath.Join(dirPath, fi.Name()))
	}
	return nil
}

func openTargetZipFile(fileName string, flags int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(fileName+".zip", flags, perm)
}

// root is the root directory name
// TODO: return an error instead of Fatal
func writeZipFiles(files []string, zipWriter *zip.Writer, root string) {
	for _, filepath := range files {
		filepathRel, err := getPathRelativeToRoot(filepath, root)
		if err != nil {
			log.Fatal("cannot get path relative to root")
		}
		if err := appendFileToZipWriter(filepath, filepathRel, zipWriter); err != nil {
			log.Fatalf("Failed to add file %s to zip: %s", filepath, err)
		}
	}
}

// getPathRelativeToRoot turns a path relative to the working directory or an
// absolute path into a path relative to root.
// For example, with root=="dst", "dst/childDir/file.txt" and
// "/User/home/dst/childDir/file.txt" will produce "childDir/file.txt"
func getPathRelativeToRoot(filePath, root string) (string, error) {
	if len(root) == 0 {
		return "", fmt.Errorf("root cannot be empty")
	}

	_, after, found := strings.Cut(filePath, root)
	if !found {
		return "", fmt.Errorf("filepath does not contain root")
	}
	return strings.TrimPrefix(after, string(os.PathSeparator)), nil
}

// filepath is a path relative to the working directory or is an absolute path
// filepathRel is a path relative to dst
func appendFileToZipWriter(filepath, filepathRel string, zipw *zip.Writer) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", filepath, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filepathRel)
	if err != nil {
		msg := "Failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filepath, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Failed to write %s to zip: %s", filepath, err)
	}

	return nil
}
