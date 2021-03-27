package tests

import (
	"archive/zip"
	"bytes"
	. "github.com/bukowa/wdgo/src/git"
	. "github.com/bukowa/wdgo/src/zip"
	"os/exec"
	"testing"
)

func Test_GitZip(t *testing.T) {
	r, err := NewTempRepository("https://github.com/bukowa/wdgo.git")
	if err != nil {
		t.Error(err)
	}
	cmdErr(r.Init(), t)
	cmdErr(r.RemoteAddOrigin(), t)
	cmdErr(r.Cmd("pull", "origin", "HEAD"), t)

	w := bytes.NewBuffer(nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	zw := zip.NewWriter(w)
	err = NewFilePathWalker(r.JoinAbs("src"), "/", zw, 0660).WalkDir(nil)
	if err != nil {
		t.Error(err)
	}
	zw.Close()
}

func cmdErr(cmd *exec.Cmd, t *testing.T) {
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Error(string(b))
		t.FailNow()
	}
}
