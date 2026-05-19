package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	PathSeparator = string(os.PathSeparator)
)

type PathInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

func (p PathInfo) Name() string {
	return p.name
}

func (p PathInfo) Size() int64 {
	return p.size
}

func (p PathInfo) Mode() fs.FileMode {
	return p.mode
}

func (p PathInfo) ModTime() time.Time {
	return p.modTime
}

func (p PathInfo) IsDir() bool {
	return p.isDir
}

func (p PathInfo) Sys() any {
	return p.sys
}

func ParsePath(path string) os.FileInfo {
	pInfo := PathInfo{
		name:  path,
		isDir: filepath.Ext(path) == "",
	}

	if info, err := Info(path); err != nil && info != nil {
		pInfo.size = info.Size()
		pInfo.mode = info.Mode()
		pInfo.modTime = info.ModTime()
		pInfo.isDir = info.IsDir()
		pInfo.sys = info.Sys()
	}

	return pInfo
}
