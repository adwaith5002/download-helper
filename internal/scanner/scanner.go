package scanner

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/adwaith5002/download-helper/pkg/fileinfo"
)

func Scan(root string) ([]fileinfo.FileInfo, error) {
	var files []fileinfo.FileInfo
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), "~$") {
			return nil
		}
		if d.Name() == "desktop.ini" {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		files = append(files, fileinfo.FileInfo{
			Path:      path,
			Name:      filepath.Base(path),
			Extension: filepath.Ext(path),
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			Category:  categorize(filepath.Ext(path)),
		})
		return nil
	})
	return files, err
}
