package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func interactiveMode() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("--- Gallery Generator ---")
	fmt.Println("Interactive Mode")
	fmt.Println()

	fmt.Print("Enter image folder path: ")
	scanner.Scan()
	folder := strings.TrimSpace(scanner.Text())
	if folder == "" {
		fmt.Println("Error: Folder path is required.")
		waitForKey()
		return
	}

	fmt.Print("Gallery Title [Photos]: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		title = "Photos"
	}

	fmt.Print("Output Filename [gallery.html]: ")
	scanner.Scan()
	output := strings.TrimSpace(scanner.Text())
	if output == "" {
		output = "gallery.html"
	}

	fmt.Print("Create Zip archive? (y/n) [y]: ")
	scanner.Scan()
	zipInput := strings.ToLower(strings.TrimSpace(scanner.Text()))
	noZip := zipInput == "n" || zipInput == "no"

	fmt.Println()
	runGenerator(folder, title, output, noZip)

	waitForKey()
}

func waitForKey() {
	fmt.Println()
	fmt.Print("Press 'Enter' to exit...")
	bufio.NewScanner(os.Stdin).Scan()
}

func main() {
	// Define flags
	titlePtr := flag.String("title", "Photos", "Title of the gallery")
	outputPtr := flag.String("output", "gallery.html", "Output HTML filename")
	noZipPtr := flag.Bool("no-zip", false, "Skip creating a zip archive")

	flag.Parse()

	// Check for positional argument (image folder)
	args := flag.Args()
	if len(args) < 1 {
		// No arguments -> Interactive Mode
		interactiveMode()
		return
	}

	// CLI Mode
	folder := args[0]
	runGenerator(folder, *titlePtr, *outputPtr, *noZipPtr)
}
