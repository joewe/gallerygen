package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// PageData holds data for the HTML template
type PageData struct {
	Title      string
	Images     []string
	FolderName string
	ZipName    string
}

func runGenerator(folder, title, outputFilename string, noZip bool) {
	// validate folder
	info, err := os.Stat(folder)
	if err != nil || !info.IsDir() {
		log.Printf("Error: Path '%s' is not a valid directory: %v", folder, err)
		return
	}

	// Get absolute path of folder and determine output directory
	absFolder, err := filepath.Abs(folder)
	if err != nil {
		log.Printf("Error to get absolute path: %v", err)
		return
	}

	// Output directory is the parent directory of the image folder
	outputDir := filepath.Dir(absFolder)
	folderName := filepath.Base(absFolder)

	// Create thumbnails directory
	thumbDir := filepath.Join(outputDir, "thumbnails")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		log.Printf("Error creating thumbnails directory: %v", err)
		return
	}

	var images []string
	// Walk through folder for image files
	entries, err := os.ReadDir(folder)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		return
	}

	fmt.Println("Scanning images and generating thumbnails...")
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		lower := strings.ToLower(e.Name())
		if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") ||
			strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".webp") ||
			strings.HasSuffix(lower, ".gif") {

			images = append(images, e.Name())

			// Generate thumbnail
			srcPath := filepath.Join(absFolder, e.Name())
			destPath := filepath.Join(thumbDir, e.Name())

			// Check if thumbnail already exists to avoid re-generating
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				fmt.Printf("Generating thumbnail for %s...\n", e.Name())
				if err := generateThumbnail(srcPath, destPath); err != nil {
					fmt.Printf("Warning: Failed to generate thumbnail for %s: %v\n", e.Name(), err)
				}
			}
		}
	}

	if len(images) == 0 {
		fmt.Println("No images found in the specified directory.")
		return
	}

	zipName := ""
	if !noZip {
		zipName = "photos.zip"
		zipPath := filepath.Join(outputDir, zipName)
		fmt.Println("Zipping folder...")
		err := zipFolder(absFolder, zipPath)
		if err != nil {
			fmt.Printf("Error creating zip: %v\n", err)
			zipName = "" // Don't show link if zip failed
		} else {
			fmt.Println("✓ Successfully created photos.zip at:", zipPath)
		}
	}

	// Prepare data for template
	data := PageData{
		Title:      title,
		Images:     images,
		FolderName: folderName,
		ZipName:    zipName,
	}

	// Parse and execute template
	tmpl, err := template.New("gallery").Parse(htmlTemplateStr)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return
	}

	outputPath := filepath.Join(outputDir, outputFilename)
	f, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating output file: %v", err)
		return
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		log.Printf("Error executing template: %v", err)
		return
	}

	fmt.Println("Ready! HTML generated:", outputPath)
}
