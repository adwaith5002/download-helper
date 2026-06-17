package organizer

import (
	"path/filepath"
	"os"
	"github.com/adwaith5002/download-helper/pkg/fileinfo"
)

type Plan struct {
	From   string
	To     string
	IsDupe bool
}

var categoryFolders = map[fileinfo.Category]string{
	fileinfo.Image:      "Images",
	fileinfo.Document:   "Documents",
	fileinfo.Video:      "Video",
	fileinfo.Audio:      "Audio",
	fileinfo.Archive:    "Archives",
	fileinfo.Code:       "Code",
	fileinfo.Executable: "Executables",
}

func isAlreadyOrganized(path string) bool {
	parent := filepath.Base(filepath.Dir(path))
	for _, folder := range categoryFolders {
		if parent == folder {
			return true
		}
	}
	return parent == "Duplicates"
}

func BuildPlan(files []fileinfo.FileInfo, dupes [][]fileinfo.FileInfo, root string) []Plan {
	dupeSet := make(map[string]bool)
	var plans []Plan
	for _, group := range dupes {
		for _, f := range group[1:] { // skip[0], it's the "original"
			dupeSet[f.Path] = true
		}
	}
	for _, f := range files {
		// 1. Skip Unknown category
		if f.Category == fileinfo.Unknown {
			continue
		}
		// 2. Check if already in a subfolder (filepath.Dir and filepath.Base)
		if isAlreadyOrganized(f.Path) {
			continue
		}
		var folder string
		// 3. Determine destination folder:
		if dupeSet[f.Path] { // ✅ — true means it's a duplicate
			folder = "Duplicates"
		} else {
			folder = categoryFolders[f.Category]
		}

		dest := filepath.Join(root, folder, f.Name)
		plans = append(plans, Plan{
			From:   f.Path,
			To:     dest,
			IsDupe: dupeSet[f.Path],
		})
	}
	return plans
}
func Execute(plans []Plan) error {
    for _, p := range plans {
		err := os.MkdirAll(filepath.Dir(p.To), 0755)
		if err!=nil{
			return err
		}
		err = os.Rename(p.From, p.To)
		if err!=nil{
			return err
		}
	}
    return nil
}