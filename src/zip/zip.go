package zip

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File interface {
	// Open is responsible for returning io.ReadCloser for a path.
	// Implementations working on File have to respect a case when
	// returned readCloser and err are both nil possibly meaning the path can be skipped.
	Open(path string) (readCloser io.ReadCloser, err error)
	// Close if responsible for closing readCloser returned by Open.
	// It has to respect a case when both; readCloser && err are nil.
	Close(readCloser io.ReadCloser, err error) error
	// Copy is responsible for copying data into zip writer.
	// Called between Open and Close only if; Open returned non nil readCloser and nil err.
	Copy(r io.Reader, w io.Writer) (written int64, err error)
}


// FilePathDenyFunc returns true if path is denied and won't be processed.
type FilePathDenyFunc func(path string, ds fs.DirEntry, err error) (bool, error)

type FilePathWalker interface {
	File
	WalkDir(denyFunc FilePathDenyFunc) error
}

func NewFilePathWalker(srcDir, dstDir string, w *zip.Writer, perm os.FileMode) FilePathWalker {
	return &filePathWalker{
		srcDir:   srcDir,
		dstDir:   dstDir,
		writer: w,
		perm: perm,
	}
}

type filePathWalker struct {
	srcDir string
	dstDir string
	writer *zip.Writer
	perm os.FileMode
}

func (f filePathWalker) Open(path string) (readCloser io.ReadCloser, err error) {
	return os.OpenFile(path, os.O_RDONLY, f.perm)
}

func (f filePathWalker) Close(readCloser io.ReadCloser, err error) error {
	if readCloser != nil {
		return readCloser.Close()
	}
	return nil
}

func (f filePathWalker) Copy(r io.Reader, w io.Writer) (written int64, err error) {
	return io.Copy(w, r)
}

func (f filePathWalker) WalkDir(denyFunc FilePathDenyFunc) error {
	return filepath.WalkDir(f.srcDir, func(path string, ds fs.DirEntry, err error) error {
		// path: /var/www/html/wp-content/plugins/akimset/index.html
		// srcPath: /var/www/html/wp-content/plugins
		// dstPath: /plugins

		// check for denied path
		if denyFunc != nil {
			if denied, err := denyFunc(path, ds, err); err != nil || denied {
				return err
			}
		}

		// open file to be copied into zip archive
		r, err := f.Open(path)
		if err != nil || r == nil {
			return err
		}

		// path as seen in zip archive
		zipPath := filepath.ToSlash(filepath.Join(filepath.Clean(f.dstDir), filepath.Clean(path[len(filepath.Clean(f.srcDir)):])))

		// handle dir
		if ds.IsDir() {

			// zip implementation expects "/" if it's a dir
			if !strings.HasPrefix(zipPath, "/") {
				zipPath += "/"
			}

			// create dir
			_, err = f.writer.Create(zipPath)
			if err != nil {
				return err
			}
			return err
		}

		// create regular file
		w, err := f.writer.Create(zipPath)
		if err != nil {
			return err
		}

		// copy file content
		_, err = f.Copy(r, w)
		return err
	})
}
