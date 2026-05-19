package files

import "os"

const (
	DefaultFileMode = 0600
)

// Write writes content to the file at filePath
// will create the file at filePath with perm permissions if the file does not exist
// perm defaults to DefaultFileMode if not provided
func Write[T ~string](filePath string, content T, perm ...os.FileMode) error {
	var mode os.FileMode = DefaultFileMode

	if len(perm) > 0 {
		mode = perm[0]
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)

	if err != nil {
		return err
	}

	defer f.Close()

	bytes := []byte(content)
	total := len(bytes)

	var written int

	for {
		i, e := f.Write(bytes)

		if e != nil {
			return e
		}

		if i == total {
			break
		}

		written += i
		bytes = bytes[i:]
	}

	return nil
}

// MakeDir creates a directory named path, along with any necessary parents
func MakeDir(path string, perm ...os.FileMode) error {
	var mode os.FileMode = DefaultFileMode

	if len(perm) > 0 {
		mode = perm[0]
	}

	return os.MkdirAll(path, mode)
}

func Touch(filePath string, perm ...os.FileMode) error {
	return Write(filePath, "")
}

func Exists(filePath string) bool {
	_, err := os.Stat(filePath)

	return err == nil
}

func Read(filePath string) (string, error) {
	c, err := os.ReadFile(filePath)

	return string(c), err
}

func Info(filePath string) (os.FileInfo, error) {
	return os.Stat(filePath)
}
