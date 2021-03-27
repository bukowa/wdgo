package tests

import (
	"archive/zip"
	. "github.com/bukowa/wdgo/src/git"
	. "github.com/bukowa/wdgo/src/zip"
	"os"
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

	f, err := os.OpenFile("x.zip", os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	w := zip.NewWriter(f)
	err = NewFilePathWalker(r.JoinAbs("src"), "/", w, 0660).WalkDir(nil)
	if err != nil {
		t.Error(err)
	}
	w.Close()
	f.Close()
}

func cmdErr(cmd *exec.Cmd, t *testing.T) {
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Error(string(b))
		t.FailNow()
	}
}
