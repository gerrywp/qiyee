package configs

import (
	"os"
	"path/filepath"
)

// getRootPath return current workspace dir
func GetRootPath() string {
	d, _ := os.Getwd()
	return filepath.Dir(d)
}
