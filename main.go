package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// PageData holds data for the HTML template
type PageData struct {
	Title      string
	Images     []string
	FolderName string
	ZipName    string
}

// HTML template with Lightbox
const htmlTemplateStr = `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; color: #333; }
        h1 { text-align: center; margin-bottom: 30px; }
        .download-link { display: block; text-align: center; margin-bottom: 40px; }
        .download-link a { background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; transition: background 0.2s; }
        .download-link a:hover { background: #0056b3; }
        
        .gallery { 
            column-width: 300px;
            column-gap: 15px;
            max-width: 1200px; 
            margin: 0 auto; 
        }
        .gallery-item { 
            break-inside: avoid; 
            margin-bottom: 15px; 
        }
        .gallery img { 
            width: 100%; 
            height: auto; 
            border-radius: 8px; 
            display: block; 
            cursor: zoom-in; 
            box-shadow: 0 2px 5px rgba(0,0,0,0.1); 
            transition: transform 0.2s;
        }
        .gallery img:hover { transform: translateY(-2px); box-shadow: 0 4px 8px rgba(0,0,0,0.15); }

        /* Lightbox */
        .lightbox {
            display: none;
            position: fixed;
            z-index: 1000;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.9);
            justify-content: center;
            align-items: center;
        }
        .lightbox.active { display: flex; }
        .lightbox img {
            max-width: 90%;
            max-height: 90vh;
            border-radius: 4px;
            box-shadow: 0 0 20px rgba(0,0,0,0.5);
        }
        .lightbox-close {
            position: absolute;
            top: 20px;
            right: 30px;
            color: white;
            font-size: 40px;
            cursor: pointer;
            user-select: none;
        }
        .lightbox-nav {
            position: absolute;
            top: 50%;
            transform: translateY(-50%);
            color: white;
            font-size: 50px;
            cursor: pointer;
            user-select: none;
            padding: 20px;
            background: rgba(0,0,0,0.2);
            border-radius: 5px;
            transition: background 0.2s;
        }
        .lightbox-nav:hover { background: rgba(0,0,0,0.5); }
        .lightbox-prev { left: 20px; }
        .lightbox-next { right: 20px; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{if .ZipName}}
    <div class="download-link">
        <a href="{{.ZipName}}">Download all photos as ZIP</a>
    </div>
    {{end}}
    
    <section class="gallery">
        {{range .Images}}
        <div class="gallery-item">
            <img class="zoomable" src="thumbnails/{{.}}" data-full="{{$.FolderName}}/{{.}}" alt="{{.}}" loading="lazy">
        </div>
        {{end}}
    </section>

    <div id="lightbox" class="lightbox">
        <span class="lightbox-close">&times;</span>
        <a class="lightbox-nav lightbox-prev">&#10094;</a>
        <a class="lightbox-nav lightbox-next">&#10095;</a>
        <img id="lightbox-img" src="" alt="Lightbox Image">
    </div>

    <script>
        const lightbox = document.getElementById('lightbox');
        const lightboxImg = document.getElementById('lightbox-img');
        const closeBtn = document.querySelector('.lightbox-close');
        const prevBtn = document.querySelector('.lightbox-prev');
        const nextBtn = document.querySelector('.lightbox-next');
        
        let images = [];
        let currentIndex = 0;

        // Collect all zoomable images
        document.querySelectorAll(".zoomable").forEach((img, index) => {
            images.push(img.dataset.full);
            img.addEventListener("click", () => {
                currentIndex = index;
                showImage(currentIndex);
                lightbox.classList.add('active');
            });
        });

        function showImage(index) {
            if (index < 0) index = images.length - 1;
            if (index >= images.length) index = 0;
            currentIndex = index;
            lightboxImg.src = images[currentIndex];
        }

        closeBtn.addEventListener('click', () => {
            lightbox.classList.remove('active');
        });

        lightbox.addEventListener('click', (e) => {
            if (e.target === lightbox) {
                lightbox.classList.remove('active');
            }
        });

        prevBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            showImage(currentIndex - 1);
        });

        nextBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            showImage(currentIndex + 1);
        });

        document.addEventListener('keydown', (e) => {
            if (!lightbox.classList.contains('active')) return;
            
            if (e.key === 'Escape') {
                lightbox.classList.remove('active');
            } else if (e.key === 'ArrowLeft') {
                showImage(currentIndex - 1);
            } else if (e.key === 'ArrowRight') {
                showImage(currentIndex + 1);
            }
        });
    </script>
</body>
</html>`

// zipFolder zips the contents of a folder using native Go archive/zip
func zipFolder(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// generateThumbnail creates a resized version of the image
func generateThumbnail(srcPath, destPath string) error {
	// Open the image
	src, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}

	// Resize the image to width = 400px, preserving the aspect ratio
	dst := imaging.Resize(src, 400, 0, imaging.Lanczos)

	// Save the resulting image as JPEG
	return imaging.Save(dst, destPath)
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
