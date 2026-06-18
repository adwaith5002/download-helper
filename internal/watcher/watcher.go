package watcher

import (
    "fmt"
    "path/filepath"
    "time"

    "github.com/fsnotify/fsnotify"
    "github.com/adwaith5002/download-helper/internal/scanner"
)

func Watch(root string) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }
    defer watcher.Close()

    err = watcher.Add(root)
    if err != nil {
        return err
    }

    fmt.Println("Watching", root, "for new files...")

    for {
        select {
        case event := <-watcher.Events:
            // YOUR CODE HERE
			if event.Op != fsnotify.Create{
				continue
			}
			ext := filepath.Ext(event.Name)
			if(ext== ".tmp" || ext== ".crdownload" || ext== ".part" ){
				continue
			}
			time.Sleep(1500 * time.Millisecond)
			category:=scanner.Categorize(filepath.Ext(event.Name))
			fmt.Println("New file:",filepath.Base(event.Name), "→ category:", category)
        
		case err := <-watcher.Errors:
            fmt.Println("watch error:", err)
        }
    }
}