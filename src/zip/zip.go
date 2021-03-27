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
	Open(path string) (rc io.ReadCloser, err error)
	// Copy is responsible for copying data into zip writer.
	Copy(r io.Reader, w io.Writer) (written int64, err error)
}

// FilePathDenyFunc returns true if path is denied and won't be processed.
type FilePathDenyFunc func(path string, ds fs.DirEntry, err error) (bool, error)

type FilePathWalker interface {
	File
	WalkDir(denyFunc FilePathDenyFunc) error
}

func NewFilePathWalker(srcDir, dstDir string, w *zip.Writer, perm os.FileMode) FilePathWalker {
	d := filepath.ToSlash(dstDir)
	if strings.HasPrefix(d, "/") {
		dstDir = dstDir[1:]
	}
	if strings.HasPrefix(d, "./") {
		dstDir = dstDir[2:]
	}
	return &filePathWalker{
		srcDir: srcDir,
		dstDir: dstDir,
		writer: w,
		perm:   perm,
	}
}

type filePathWalker struct {
	srcDir string
	dstDir string
	writer *zip.Writer
	perm   os.FileMode
}

func (f filePathWalker) Open(path string) (readCloser io.ReadCloser, err error) {
	return os.OpenFile(path, os.O_RDONLY, f.perm)
}

func (f filePathWalker) Copy(r io.Reader, w io.Writer) (written int64, err error) {
	return io.Copy(w, r)
}

func (f filePathWalker) WalkDir(denyFunc FilePathDenyFunc) error {

	var newFileHeader = func(path string, ds fs.DirEntry) (*zip.FileHeader, error) {
		fi, err := ds.Info()
		if err != nil {
			return nil, err
		}
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return nil, err
		}

		// path as seen in zip archive
		fh.Name = filepath.ToSlash(filepath.Join(filepath.Clean(f.dstDir), filepath.Clean(path[len(filepath.Clean(f.srcDir)):])))

		if ds.IsDir() {
			// zip implementation expects "/" if it's a dir
			if !strings.HasPrefix(fh.Name, "/") {
				fh.Name += "/"
			}
		}
		return fh, nil
	}

	return filepath.WalkDir(f.srcDir, func(path string, ds fs.DirEntry, err error) error {

		// check for denied path
		if denyFunc != nil {
			if denied, err := denyFunc(path, ds, err); err != nil || denied {
				return err
			}
		}

		// open file to be copied into zip archive
		r, err := f.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		// create file header
		fh, err := newFileHeader(path, ds)
		if err != nil {
			return err
		}

		// create zip header
		w, err := f.writer.CreateHeader(fh)
		if err != nil {
			return err
		}

		// handle dir
		if ds.IsDir() {
			return nil
		}

		// copy file content
		_, err = f.Copy(r, w)
		return err
	})
}
