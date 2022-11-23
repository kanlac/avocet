package main

import (
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	os.Args = []string{"avocet", "create", "testdir"}
	main()
}

func TestGetFilesInDirectory(t *testing.T) {
	// var files []string
	// if err := getFilesInDirectory("testdir", &files); err != nil {
	// 	t.Error(err)
	// }
	// t.Logf("files: %+v", files[:cap(files)])
}

func TestGetPathRelativeToRoot(t *testing.T) {
	testCases := []struct {
		Path     string
		Root     string
		Expected string
	}{
		{
			"dst/childDir/file.txt",
			"dst",
			"childDir/file.txt",
		},
		{
			"/User/home/dst/childDir/file.txt",
			"dst",
			"childDir/file.txt",
		},
		{
			"foo/dst/childDir/file.txt",
			"dst",
			"childDir/file.txt",
		},
	}

	for _, tc := range testCases {
		got, err := getPathRelativeToRoot(tc.Path, tc.Root)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.Expected {
			t.Fatalf("path: %s, root: %s, expected: %s, got: %s", tc.Path, tc.Root, tc.Expected, got)
		}
	}

}
