package zip_test

import (
	"archive/zip"
	"bytes"
	. "github.com/bukowa/wdgo/src/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_Zip(t *testing.T) {
	var testPath = testWorkDir()
	log.Print(testPath)

	var testDataWant = map[string][]byte{
		"zipdest/1":   []byte("11"),
		"zipdest/3/4": []byte("44"),
		"zipdest/2":   []byte("22"),
	}

	wd, err := filepath.Abs(testPath)
	if err != nil {
		t.Error(err)
	}

	fname := "TestWorkDir_Zip.zip.zip"
	zipPath := filepath.Join(wd, "testdata", fname)
	//defer os.RemoveAll(zipPath)

	openZipFile := func(write bool) *os.File {
		var perms int
		switch write {
		case true:
			os.RemoveAll(zipPath)
			perms = os.O_CREATE | os.O_EXCL
		case false:
			perms = os.O_RDONLY
		}
		f, err := os.OpenFile(zipPath, perms, 0660)
		if err != nil {
			t.Error(err)
		}
		return f
	}

	f := openZipFile(true)

	var zw = zip.NewWriter(f)

	zipDir := filepath.Join(testPath, "testdata", "TestWorkDir_Zip", "zip")
	err = Zip(zipDir, "zipdest", zw)
	if err != nil {
		t.Error(err)
		return
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
