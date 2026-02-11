package cmd

import "fmt"

// Version information set via ldflags during build
const (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func PrintVersion() {
	fmt.Printf("monkey %s\n", Version)
	fmt.Printf("  commit: %s\n", Commit)
	fmt.Printf("  built:  %s\n", Date)
}
