package zipfs

import (
	"io/fs"
	"time"
)

// dirEntry implement fs.DirEntry, which present of zip file members list
type dirEntry struct {
	fs.FileInfo
	path string
}

func (d dirEntry) IsDir() bool {
	return d.FileInfo.IsDir()
}

func (d dirEntry) Type() fs.FileMode {
	return d.FileInfo.Mode().Type()
}

func (d dirEntry) Info() (fs.FileInfo, error) {
	return d.FileInfo, nil
}

func (d dirEntry) Name() string {
	return d.path
}

// zipInfo presentation of virtual zip root dir directory, which implement fs.FileInfo
type zipInfo struct {
	name string
}

func (d *zipInfo) Name() string {
	return d.name
}

func (d *zipInfo) IsDir() bool {
	return true
}

func (d *zipInfo) Type() fs.FileMode {
	return fs.ModeDir.Perm().Type()
}

func (d *zipInfo) ModTime() time.Time {
	return time.Now()
}

func (d *zipInfo) Mode() fs.FileMode {
	return 0755
}

func (d *zipInfo) Sys() any {
	return nil
}

func (d *zipInfo) Size() int64 {
	return 0
}
