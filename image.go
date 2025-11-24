package main

import (
	"github.com/disintegration/imaging"
)

// generateThumbnail creates a resized version of the image
func generateThumbnail(srcPath, destPath string) error {
	// Open the image
	src, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}

	// Resize the image to width = 1000px, preserving the aspect ratio
	dst := imaging.Resize(src, 1000, 0, imaging.Lanczos)

	// Save the resulting image as JPEG
	return imaging.Save(dst, destPath)
}
