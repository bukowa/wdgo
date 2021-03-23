package zip

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ZipDir creates a zip archive copying files from path in WorkDir.
// Example:
// 	zipDir := wd.JoinAbs(workingDir, "testdata", "TestWorkDir_Zip", "zip")
//	err = Zip(ZipWalkDirFunc(zipDir, "zipdest", zw), zipDir)
func Zip(src, dst string, w *zip.Writer) error {
	return filepath.WalkDir(src, ZipWalkDirFunc(src, dst, w))
}

// ZipWalkDirFunc returns fs.WalkDirFunc for zipping files.
// Files are copied from dst to src in w zip.Writer.
var ZipWalkDirFunc = func(src, dst string, w *zip.Writer) func(string, fs.DirEntry, error) error {
	return func(path string, ds fs.DirEntry, err error) error {

		// path: /var/www/html/wp-content/plugins/akimset/index.html
		// src: /var/www/html/wp-content/plugins
		// dst: /plugins
		if err != nil {
			return err
		}

		// path as seen in zip archive
		zp := filepath.ToSlash(filepath.Join(filepath.Clean(dst), filepath.Clean(path[len(filepath.Clean(src)):])))

		// handle dir
		if ds.IsDir() {

			// zip implementation expects "/" if it's a dir
			if !strings.HasPrefix(zp, "/") {
				zp += "/"
			}

			// create dir
			_, err = w.Create(zp)
			if err != nil {
				return err
			}
			return err
		}

		// create regular file
		w, err := w.Create(zp)
		if err != nil {
			return err
		}

		// open file to copy
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		defer f.Close()

		// copy file content
		_, err = io.Copy(w, f)
		return err
	}
}

