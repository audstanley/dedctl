package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// defaultImageFiles lists the image filenames to ensure exist in imgDir.
var defaultImageFiles = []string{"main_cuttle.png", "cuttle_icon.png"}

// sourceImageDir returns the path to the source images directory (configs/img/) relative to the binary.
func sourceImageDir() string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	current := dir
	for {
		imgPath := filepath.Join(current, "configs", "img")
		if _, err := os.Stat(imgPath); err == nil {
			return imgPath
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return ""
}

// ensureDefaultImages copies main_cuttle.png and cuttle_icon.png into imgDir if they're not already there.
func ensureDefaultImages(imgDir string) error {
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		return err
	}

	srcDir := sourceImageDir()
	if srcDir == "" {
		return fmt.Errorf("could not locate source images directory")
	}

	for _, name := range defaultImageFiles {
		srcPath := filepath.Join(srcDir, name)
		destPath := filepath.Join(imgDir, name)
		if _, err := os.Stat(destPath); err == nil {
			continue
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s from %s: %v\n", name, srcPath, err)
			continue
		}
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to copy %s to %s: %v\n", name, destPath, err)
		} else {
			fmt.Printf("Copied default image: %s\n", name)
		}
	}
	return nil
}
