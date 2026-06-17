package cli

import (
	"os"
	"path/filepath"
)

type Executable interface {
	GetName() string
	GetFullPath() string
}

type executable struct {
	Name     string `json:"name,omitempty"`
	FullPath string `json:"fullPath,omitempty"`
}

func (e *executable) GetName() string {
	return e.Name
}

func (e *executable) GetFullPath() string {
	return e.FullPath
}

var (
	exe *executable
)

// GetExecutable gets information about the executing binary
func GetExecutable() Executable {
	if exe != nil {
		return exe
	}

	path, pErr := os.Executable()

	if pErr != nil {
		return nil
	}

	exe = &executable{
		Name:     filepath.Base(path),
		FullPath: path,
	}

	return exe
}
