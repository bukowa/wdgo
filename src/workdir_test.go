package src_test

import (
	. "github.com/bukowa/wdgo/src"
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

func testWorkDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}
