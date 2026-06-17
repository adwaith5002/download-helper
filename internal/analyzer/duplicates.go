package analyzer

import (
	"github.com/adwaith5002/download-helper/internal/hasher"
	"github.com/adwaith5002/download-helper/pkg/fileinfo"
)

func FindDuplicates(files []fileinfo.FileInfo) ([][]fileinfo.FileInfo, error) {
	// Step 1: group by size
	bySize := make(map[int64][]fileinfo.FileInfo)
	for _, f := range files {
		bySize[f.Size] = append(bySize[f.Size], f)
	}

	// Step 2: for size groups > 1, hash and group by hash
	byHash := make(map[string][]fileinfo.FileInfo)
	for _, group := range bySize {
		if len(group) < 2 {
			continue // skip unique sizes
		}
		for _, f := range group {
			hash, err := hasher.HashFile(f.Path)
			if err != nil {
				return nil, err
			}
			byHash[hash] = append(byHash[hash], f)

		}
		// hash each file in the group
		// append to byHash[hash]
	}

	// Step 3: collect hash groups > 1 into result
	var duplicates [][]fileinfo.FileInfo
	for _, group := range byHash {
		if len(group) > 1 {
			duplicates = append(duplicates, group) // this is a duplicate group
		}
	}

	return duplicates, nil
}
