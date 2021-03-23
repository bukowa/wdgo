package src

import (
	"path/filepath"
)

// CleanBasePath returns path cleaned of basedir, as seen in the basedir root.
// Always returns a path with "/" as separator.
// Example:
//	CleanBasePath("/var/www/html", "/var/www/html/root/index.html") == "/root/index.html"
func CleanBasePath(basedir, path string) string {
	// patch as is seen in the basedir
	p := filepath.Clean(path[len(filepath.Clean(basedir)):])
	// slashes are OS dependant, here we convert them to "/"
	p = filepath.ToSlash(p)
	return p
}
