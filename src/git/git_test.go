package git

import (
	"os/exec"
	"testing"
)

func cmdErr(cmd *exec.Cmd, t *testing.T) {
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Error(string(b))
		t.FailNow()
	}
}

func TestRepository_RemoteAddOrigin(t *testing.T) {
	r, err := NewTempRepository("https://github.com/bukowa/wdgo.git")
	if err != nil {
		t.Error(err)
	}
	cmdErr(r.Init(), t)
	cmdErr(r.RemoteAddOrigin(), t)
	cmdErr(r.Cmd("pull", "origin", "HEAD"), t)
	_, err = r.WorkDir.Open("README.md")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
