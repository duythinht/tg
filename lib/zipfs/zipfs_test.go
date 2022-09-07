package zipfs_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/fs"
	"testing"

	"github.com/duythinht/tg/lib/zipfs"
)

func TestZipFSWalk(t *testing.T) {
	b := bytes.NewBuffer(nil)
	zipWriter := zip.NewWriter(b)
	w, err := zipWriter.Create("cmd/xyz/main.go")
	if err != nil {
		t.Fatalf("create zip writer failed %s", err)
	}

	fmt.Fprintf(w, `
	package main
	func main() {
		fmt.Println("Hello, world!")
	}
	`)

	zipWriter.Close()

	zipData := b.Bytes()

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))

	if err != nil {
		t.Fatalf("create zip reader %s", err)
	}

	fsys := zipfs.NewFS(zipReader)

	fs.WalkDir(fsys, "", func(path string, d fs.DirEntry, err error) error {

		// skip dir
		if path == "" {
			if !d.IsDir() {
				t.Logf("root path `` should be dir")
				t.Fail()
			}
			return nil
		}

		if path != "cmd/xyz/main.go" {
			t.Logf("path `%s` is wrong", path)
			t.Fail()
		}
		return nil
	})
}
