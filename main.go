package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()

	// Use the first positional argument as the root directory; default to "." if not provided.
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	// Tilde Expansion
	// Expand tilde if present (e.g., "~/D" -> "/home/username/D")
	if strings.HasPrefix(root, "~") {
		home, err := os.UserHomeDir()  // get the current user's home
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

	// Walk through the directory tree and count files and directories.
	walkDir(root, &totalFiles, &totalDirs)

	// Print total counts.
	fmt.Printf("\nFrom: %s\n", root)
	fmt.Printf("Files: %d\n", totalFiles)
	fmt.Printf("Dirs : %d\n", totalDirs)
}

// walkDir recursively processes the given directory, counting files and directories.
func walkDir(dir string, totalFiles *int, totalDirs *int) {
	// List the directory entries.
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("failed to read directory %s: %v\n", dir, err)
		return
	}

	// Increment directory count.
	*totalDirs++

	// Iterate over each entry.
	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			// Recursively process subdirectories.
			walkDir(fullPath, totalFiles, totalDirs)
		} else {
			// Increment file count.
			*totalFiles++
		}
	}
}
