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
	var currentConcurrent, maxConcurrent int

	// Start traversing the directory tree concurrently.
	wg.Add(1)
	go walkDir(root, &totalFiles, &totalDirs, &wg, &mu, &currentConcurrent, &maxConcurrent)
	wg.Wait()

	// Print total counts and the maximum concurrent calls.
	fmt.Printf("\nFrom : %s\n", root)
	fmt.Printf("Files: %d\n", totalFiles)
	fmt.Printf("Dirs : %d\n", totalDirs)
	fmt.Printf("Max Concurrent walkDir calls: %d\n", maxConcurrent)
}

// walkDir recursively processes the given directory concurrently,
// counting files and directories and tracking concurrent calls.
func walkDir(dir string, totalFiles *int, totalDirs *int, wg *sync.WaitGroup, mu *sync.Mutex, currentConcurrent *int, maxConcurrent *int) {
	// Update concurrent call counters.
	mu.Lock()
	*currentConcurrent++
	if *currentConcurrent > *maxConcurrent {
		*maxConcurrent = *currentConcurrent
	}
	mu.Unlock()

	// Decrement current concurrency and mark this goroutine as done when function returns.
	defer func() {
		mu.Lock()
		*currentConcurrent--
		mu.Unlock()
		wg.Done()
	}()

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
			wg.Add(1)
			go walkDir(fullPath, totalFiles, totalDirs, wg, mu, currentConcurrent, maxConcurrent)
		} else {
			mu.Lock()
			*totalFiles++
			mu.Unlock()
		}
	}
}
