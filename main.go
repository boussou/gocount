package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// setting it to 10 has significant impact in perf
const maxConcurrent = 100

func main() {
	flag.Parse()

	// Use the first positional argument as the root directory; default to "." if not provided.
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	// Tilde Expansion: Expand "~" to the user's home directory.
	if strings.HasPrefix(root, "~") {
		home, err := os.UserHomeDir() // get the current user's home
		if err != nil {
			log.Fatalf("failed to get user home directory: %v", err)
		}
		if root == "~" {
			root = home
		} else if strings.HasPrefix(root, "~/") {
			root = filepath.Join(home, root[2:])
		}
	}

	var totalFiles, totalDirs int
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Create a buffered channel to limit concurrent goroutines.
	sem := make(chan struct{}, maxConcurrent)

	// Start traversing the directory tree concurrently.
	wg.Add(1)
	go walkDir(root, &totalFiles, &totalDirs, &wg, &mu, sem)
	wg.Wait()

	// Print total counts.
	fmt.Printf("\nFrom: %s\n", root)
	fmt.Printf("Files: %d\n", totalFiles)
	fmt.Printf("Dirs : %d\n", totalDirs)
}

// walkDir recursively processes the given directory concurrently,
// counting files and directories. It uses a semaphore (sem) to limit
// the number of concurrent invocations.
func walkDir(dir string, totalFiles *int, totalDirs *int, wg *sync.WaitGroup, mu *sync.Mutex, sem chan struct{}) {
	// Acquire a slot in the semaphore.
	sem <- struct{}{}
	// Release the slot when done.
	defer func() { <-sem }()
	defer wg.Done()

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("failed to read directory %s: %v\n", dir, err)
		return
	}

	// Increment directory count safely.
	mu.Lock()
	*totalDirs++
	mu.Unlock()

	// Process each entry.
	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			// Process subdirectory in a new goroutine.
			wg.Add(1)
			go walkDir(fullPath, totalFiles, totalDirs, wg, mu, sem)
		} else {
			// Safely increment file count.
			mu.Lock()
			*totalFiles++
			mu.Unlock()
		}
	}
}
