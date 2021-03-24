package zip_test

import (
	"archive/zip"
	"bytes"
	. "github.com/bukowa/wdgo/src/zip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)
//
//func Test_Zip(t *testing.T) {
//	var testPath = testWorkDir()
//	log.Print(testPath)
//
//	var testDataWant = map[string][]byte{
//		"zipdest/1":   []byte("11"),
//		"zipdest/3/4": []byte("44"),
//		"zipdest/2":   []byte("22"),
//	}
//
//	wd, err := filepath.Abs(testPath)
//	if err != nil {
//		t.Error(err)
//	}
//
//	fname := "TestWorkDir_Zip.zip.zip"
//	zipPath := filepath.Join(wd, "testdata", fname)
//	//defer os.RemoveAll(zipPath)
//	os.RemoveAll(zipPath)
//	f,err := os.OpenFile(zipPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0660)
//	if err != nil {
//		t.Error(err)
//	}
//	var zw = zip.NewWriter(f)
//
//	zipDir := filepath.Join(testPath, "testdata", "TestWorkDir_Zip", "zip")
//	err = Zip(zipDir, "zipdest", zw, nil)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	if err = zw.Close(); err != nil {
//		t.Error(err)
//	}
//
//	if err = f.Close(); err != nil {
//		t.Error(err)
//	}
//
//	r, err := zip.OpenReader(zipPath)
//	if err != nil {
//		t.Error(err)
//	}
//	defer r.Close()
//
//	for k, v := range testDataWant {
//		f, err := r.Open(k)
//		if err != nil {
//			t.Error(err)
//		}
//		b, err := io.ReadAll(f)
//		if err != nil {
//			t.Error(err)
//		}
//		if !bytes.Equal(v, b) {
//			t.Error(v, b)
//		}
//	}
//
//}

// testWorkDir returns path of this file.
func testWorkDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func Test_filePathWalker_WalkDir(t *testing.T) {
	type fields struct {
		srcDir string
		dstDir string
		writer *zip.Writer
		perm   os.FileMode
		want map[string][]byte
	}
	type args struct {
		denyFunc FilePathDenyFunc
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    ".",
			fields:  fields{
				srcDir: filepath.Join(testWorkDir(), "testdata", "TestNewFilePathWalker", "zip"),
				dstDir: "zipdest",
				// set later in test
				writer: nil,
				perm:   0600,
				want: map[string][]byte{
					"zipdest/1":   []byte("11"),
					"zipdest/3/4": []byte("44"),
					"zipdest/2":   []byte("22"),
				},
			},
			args:    args{
				denyFunc: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			tt.fields.writer = zip.NewWriter(w)
			f := NewFilePathWalker(
				tt.fields.srcDir,
				tt.fields.dstDir,
				tt.fields.writer,
				tt.fields.perm)
			if err := f.WalkDir(tt.args.denyFunc); (err != nil) != tt.wantErr {
				t.Errorf("WalkDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.fields.writer.Close(); err != nil {
				t.Error(err)
			}

			zipBytes, err := io.ReadAll(w)
			if err != nil {
				t.Error(err)
			}

			r := bytes.NewReader(zipBytes)
			zReader, err := zip.NewReader(r, r.Size())
			if err != nil {
				t.Error(err)
			}

			// compare
			for k, v := range tt.fields.want {
				f, err := zReader.Open(k)
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
		})
	}
}