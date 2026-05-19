package files

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Zip(path string) (string, error) {
	zipFile := strings.TrimSuffix(path, PathSeparator) + ".zip"

	zFile, zfErr := os.Create(zipFile)

	if zfErr != nil {
		return "", zfErr
	}

	defer zFile.Close()

	zw := zip.NewWriter(zFile)

	defer zw.Close()

	if err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, rpErr := filepath.Rel(path, filePath)

		if rpErr != nil {
			return rpErr
		}

		zf, zErr := zw.Create(relPath)

		if zErr != nil {
			return zErr
		}

		f, fErr := os.Open(filePath)

		if fErr != nil {
			return fErr
		}

		defer f.Close()

		if _, cErr := io.Copy(zf, f); cErr != nil {
			return cErr
		}

		return nil
	}); err != nil {
		return "", err
	}

	_ = os.RemoveAll(path)

	return zipFile, nil
}

func Unzip(zipFile string, dest string) error {
	z, zErr := zip.OpenReader(zipFile)

	if zErr != nil {
		return zErr
	}

	defer z.Close()

	for _, f := range z.File {
		fPath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fPath, 0750); err != nil {
				return err
			}

			continue
		}

		fDest, fdErr := os.Create(fPath)

		if fdErr != nil {
			return fdErr
		}

		fSrc, fsErr := f.Open()

		if fsErr != nil {
			_ = fDest.Close()
			return fsErr
		}

		if _, err := io.Copy(fDest, fSrc); err != nil {
			_ = fSrc.Close()
			_ = fDest.Close()
			return err
		}

		_ = fSrc.Close()
		_ = fDest.Close()
	}

	_ = os.Remove(zipFile)

	return nil
}
