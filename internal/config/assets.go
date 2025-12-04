package config

import (
	// "fmt"
	"os"
	"path/filepath"
)

const WorkingDir = "GoApps/NovelParser"

// type AssetPathS struct {
// 	Path string
// }

// func (t AssetPathS) String() string {
// 	dir, _ := os.UserHomeDir()

// 	return filepath.Join(dir, WorkingDir, t.Path)
// }

func AssetPath(path ...string) string {
	dir, _ := os.UserHomeDir()
	p := []string{dir, WorkingDir}
	p = append(p, path...)

	return filepath.Join(p...)
}
