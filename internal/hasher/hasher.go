package hasher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	return hash, err
}
