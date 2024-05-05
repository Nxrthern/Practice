package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Get the directory where the Makefile resides
	makefileDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Determine the root directory of the Git repository
	rootDir, err := gitRoot(makefileDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print the root directory
	fmt.Println(rootDir)
}

func gitRoot(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
