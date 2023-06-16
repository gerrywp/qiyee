package configs

import (
	"fmt"
	"os"
	"path/filepath"
)

// getRootPath return current workspace dir
func getRootPath() string {
	d, _ := os.Getwd()
	return filepath.Dir(d)
}

// GetAbsPath eval absolute path
func GetAbsPath(p string) string {
	fmt.Println(filepath.Join(getRootPath(), p))
	return filepath.Join(getRootPath(), p)
}
