package analyzer

import (
	"fmt"
	"time"

	"github.com/adwaith5002/download-helper/pkg/fileinfo"
)
type Recommendation struct{
	Message string
	Priority Priority
	Files []fileinfo.FileInfo
}
type Priority int

const (
	Info     Priority = iota
	Warning
	Critical
)
func (p Priority) String() string{
	switch p {
	case Info:
		return "Info"
		
	case Warning:
		return "Warning"

	case Critical:
		return "Critical"

	default:
	    return "Unknown"
	}
	
}	

func Recommend(files []fileinfo.FileInfo, duplicates [][]fileinfo.FileInfo) []Recommendation{
	var count int
	var exec int
	var recommendations []Recommendation

    // Rule 1: duplicates
    if len(duplicates) > 0 {
        var wasted int64
        var dupeFiles []fileinfo.FileInfo
        for _, group := range duplicates {
            for _, f := range group[1:] {
                wasted += f.Size
                dupeFiles = append(dupeFiles, f)
            }
        }
        recommendations = append(recommendations, Recommendation{
            Message:  fmt.Sprintf("Found %d duplicate groups wasting %d bytes", len(duplicates), wasted),
            Priority: Warning,
            Files:    dupeFiles,
        })
    }
	var timeFiles []fileinfo.FileInfo
	var execFiles []fileinfo.FileInfo
	for _, f := range files {
		if time.Since(f.ModTime)>6 * 30 * 24 * time.Hour{
			count+=1
			timeFiles = append(timeFiles, f)
		}
		if f.Category == fileinfo.Executable{
			exec+=1
			execFiles = append(execFiles, f)
		}
		
	}
	if count>5{
		recommendations = append(recommendations, Recommendation{
            Message:  fmt.Sprintf("Found %d older files spanning more than six months", count),
            Priority: Warning,
            Files:    timeFiles,
        })
	}
	if exec>3{
		recommendations = append(recommendations, Recommendation{
            Message:  fmt.Sprintf("Found %d exec files",exec),
            Priority: Warning,
            Files:    execFiles,
        })		
	}
	return recommendations
}
