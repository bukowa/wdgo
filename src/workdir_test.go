package src_test

import (
	"archive/zip"
	"bytes"
	. "github.com/bukowa/wdgo/src"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWorkDir_Abs(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileDir := filepath.ToSlash(filepath.Dir(thisFilePath))

	tests := []struct {
		name    string
		wd      string
		wantDir string
		wantErr bool
	}{
		{
			name:    ".",
			wd:      ".",
			wantDir: thisFileDir,
			wantErr: false,
		},
		{
			name:    "./",
			wd:      "./",
			wantDir: thisFileDir,
			wantErr: false,
		},
		{
			name:    "./.",
			wd:      "./.",
			wantDir: thisFileDir,
			wantErr: false,
		},
		{
			name:    "advanced",
			wd:      "advanced",
			wantDir: filepath.ToSlash(filepath.Join(thisFileDir, "advanced")),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd, _ := NewWorkDir(tt.wd)
			gotDir := wd.Abs()
			if gotDir != tt.wantDir {
				t.Errorf("Abs() gotDir = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func TestWorkDir_Zip(t *testing.T) {
	var testPath = testWorkDir()
	log.Print(testPath)

	var testDataWant = map[string][]byte{
		"zip/1":   []byte("11"),
		"zip/3/4": []byte("44"),
		"zip/2": []byte("22"),
	}

	wd, err := NewWorkDir(testPath)
	if err != nil {
		t.Error(err)
	}

	fname := "TestWorkDir_Zip.zip.zip"
	zipPath := filepath.Join("testdata", fname)
	defer wd.RemoveAll(zipPath)

	openZipFile := func(write bool) *os.File{
		var perms int
		switch write {
		case true:
			wd.RemoveAll(zipPath)
			perms = os.O_CREATE|os.O_EXCL
		case false:
			perms = os.O_RDONLY
		}
		f, err := wd.OpenFile(zipPath, perms, 0660)
		if err != nil {
			t.Error(err)
		}
		return f
	}

	f := openZipFile(true)

	var zw = zip.NewWriter(f)

	err = wd.Zip(ZipWalkDirFunc(filepath.Join(testPath, "testdata/TestWorkDir_Zip"), zw), "testdata/TestWorkDir_Zip/zip")
	if err != nil {
		t.Error(err)
	}

	if err = zw.Close(); err != nil {
		t.Error(err)
	}

	if err = f.Close(); err != nil {
		t.Error(err)
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	for k, v := range testDataWant {
		f, err := r.Open(k)
		if err != nil {
			t.Error(err)
		}
		b, err := io.ReadAll(f)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(v, b) {
			t.Error(v, b)
		}
	}

}


func testWorkDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}