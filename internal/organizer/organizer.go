package organizer

import (
	"os"
	"path/filepath"

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
        for _, f := range group[1:] {
            dupeSet[f.Path] = true
        }
    }

    for _, f := range files {
        if f.Category == fileinfo.Unknown {
            continue
        }
        if isAlreadyOrganized(f.Path) {
            continue
        }
        var folder string
        if dupeSet[f.Path] {
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
		if _, err := os.Stat(p.From); os.IsNotExist(err) {
			continue // source no longer exists, skip it
		}
		err := os.MkdirAll(filepath.Dir(p.To), 0755)
		if err != nil {
			return err
		}
		err = os.Rename(p.From, p.To)
		if err != nil {
			return err
		}
	}
	return nil
}
