package fileinfo

import "time"

type FileInfo struct {
	Path      string
	Name      string
	Size      int64
	Extension string
	ModTime   time.Time
	Category Category
}
type Category int

const (
	Unknown Category = iota
	Image
	Document
	Video
	Audio
	Archive
	Code
	Executable
)
func (c Category) String() string {
    switch c {
    case Image:
        return "Image"
    case Document:
        return "Document"
    case Video:
        return "Video"
    case Audio:
        return "Audio"
    case Archive:
        return "Archive"
    case Code:
        return "Code"
    case Executable:
        return "Executable"
    default:
        return "Unknown"
    }
}