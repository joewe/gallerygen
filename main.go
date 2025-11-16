package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Simple HTML template; {{IMAGES}} will be replaced
const htmlTemplate = `<!doctype html>
<html><head><meta charset="utf-8" /><title>Galerie</title>
<style>
body{font-family:sans-serif;background:#f5f5f5;margin:0;padding:20px}
.gallery{column-width:300px;column-gap:10px}
.gallery img{width:100%;margin-bottom:10px;border-radius:6px;display:block}
</style></head><body>
<h1>Photos</h1>
<p><a href="photos.zip">Download all photos as ZIP</a></p>
<section class="gallery">
{{IMAGES}}
</section>
<script>
document.querySelectorAll(".zoomable").forEach(img => {
  img.addEventListener("click", () => {
    window.open(img.src, "_blank");
  });
});
</script>
</body></html>`

// Zip a folder using the 'zip' command
func ZipWithZip(sourceFolder, targetZip string) error {
	// Get absolute paths
	absSource, err := filepath.Abs(sourceFolder)
	if err != nil {
		return err
	}
	absTarget, err := filepath.Abs(targetZip)
	if err != nil {
		return err
	}

	// Change to parent directory of source folder
	parentDir := filepath.Dir(absSource)
	folderName := filepath.Base(absSource)

	cmd := exec.Command("zip", "-r", absTarget, folderName)
	cmd.Dir = parentDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CheckCommandExists checks if a command is available
func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gallerygen <image-folder> <output.html>")
		return
	}

	folder := os.Args[1]
	outputFilename := os.Args[2]

	// validate folder
	info, err := os.Stat(folder)
	if err != nil || !info.IsDir() {
		log.Fatalf("Path is not a valid directory: %v", err)
	}

	// Get absolute path of folder and determine output directory
	absFolder, err := filepath.Abs(folder)
	if err != nil {
		log.Fatalf("Error to get absolute path: %v", err)
	}

	// Output directory is the parent directory of the image folder
	outputDir := filepath.Dir(absFolder)
	folderName := filepath.Base(absFolder)

	var images []string
	// Walk through folder for image files
	filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		lower := strings.ToLower(d.Name())
		if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") ||
			strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".webp") ||
			strings.HasSuffix(lower, ".gif") {
			images = append(images, d.Name())
		}
		return nil
	})

	// Build HTML img tags
	var b strings.Builder
	for _, img := range images {
		b.WriteString(fmt.Sprintf("<img class=\"zoomable\" src=\"%s/%s\" alt=\"\" loading=\"lazy\">\n", folderName, img))
	}

	// Zip the folder using 'zip' command
	zipPath := filepath.Join(outputDir, "photos.zip")
	fmt.Println("Zipping folder...")
	if CheckCommandExists("zip") {
		err := ZipWithZip(absFolder, zipPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("✓ Successfully created photos.zip at:", zipPath)
		}
	} else {
		fmt.Println("✗ 'zip' command not found")
	}

	// Replace marker
	outputPath := filepath.Join(outputDir, outputFilename)
	html := strings.Replace(htmlTemplate, "{{IMAGES}}", b.String(), 1)

	// Write to output file
	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Println("Ready! HTML generated:", outputPath)
}
