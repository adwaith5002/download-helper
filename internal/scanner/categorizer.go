package scanner

import "github.com/adwaith5002/download-helper/pkg/fileinfo"

var extMap = map[string]fileinfo.Category{
	".jpg":  fileinfo.Image,
	".jpeg": fileinfo.Image,
	".png":  fileinfo.Image,
	".gif":  fileinfo.Image,
	".webp": fileinfo.Image,
	".avif": fileinfo.Image,
	".pdf":  fileinfo.Document,
	".docx": fileinfo.Document,
	".xlsx": fileinfo.Document,
	".pptx": fileinfo.Document,
	".txt":  fileinfo.Document,
	".csv":  fileinfo.Document,
	".md":   fileinfo.Document,
	".json": fileinfo.Document,
	".mp4":  fileinfo.Video,
	".mkv":  fileinfo.Video,
	".mov":  fileinfo.Video,
	".mp3":  fileinfo.Audio,
	".wav":  fileinfo.Audio,
	".zip":  fileinfo.Archive,
	".rar":  fileinfo.Archive,
	".tar":  fileinfo.Archive,
	".iso":  fileinfo.Archive,
	".ino":  fileinfo.Code, // Arduino
	".jsx":  fileinfo.Code, // React
	".jar":  fileinfo.Code, // Java archive
	".go":   fileinfo.Code,
	".py":   fileinfo.Code,
	".js":   fileinfo.Code,
	".exe":  fileinfo.Executable,
	".msi":  fileinfo.Executable,
}

func Categorize(ext string) fileinfo.Category {
	category, ok := extMap[ext]
	if !ok {
		return fileinfo.Unknown
	}
	return category
}
