// zipfs implement io/fs abstraction to manipulate zipfile easily as file system
package zipfs

import (
	"archive/zip"
	"fmt"
	"io/fs"
)

type FS struct {
	*zip.Reader
}

func NewFS(r *zip.Reader) *FS {
	return &FS{
		Reader: r,
	}
}

// Open implement for fs.FS Open method
func (fsys *FS) Open(name string) (fs.File, error) {

	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("path `%s` invalid", name)
	}

	f, err := fsys.Reader.Open(name)

	if err != nil {
		return nil, fmt.Errorf("can not open `%s`: %w", name, err)
	}

	return f, nil
}

// ReadDir implement for fs.ReadDirFS
func (fsys *FS) ReadDir(name string) ([]fs.DirEntry, error) {

	entries := make([]fs.DirEntry, 0)

	for _, f := range fsys.File {

		if f.FileInfo().IsDir() {
			continue
		}

		entries = append(entries, dirEntry{
			FileInfo: f.FileInfo(),
			path:     f.Name,
		})

	}
	return entries, nil
}

// Stat implement for manupulate fs.StatFS
func (fsys *FS) Stat(name string) (fs.FileInfo, error) {

	if name == "" {
		return &zipInfo{}, nil
	}

	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("path `%s` invalid", name)
	}

	f, err := fsys.Reader.Open(name)

	if err != nil {
		return nil, fmt.Errorf("can not open `%s`: %w", name, err)
	}

	return f.Stat()
}
