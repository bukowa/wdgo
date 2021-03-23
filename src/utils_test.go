package src

import "testing"

func TestCleanBasePath(t *testing.T) {
	type args struct {
		basedir string
		path    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "file",
			args: args{
				basedir: "/var/www/html",
				path:    "/var/www/html/root/index.html",
			},
			want: "/root/index.html",
		},
		{
			name: "directory",
			args: args{
				basedir: "/var/www/html/",
				path:    "/var/www/html/root",
			},
			want: "/root",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanBasePath(tt.args.basedir, tt.args.path); got != tt.want {
				t.Errorf("CleanBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
