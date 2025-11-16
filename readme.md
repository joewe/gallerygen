# Gallery Generator

A simple Go tool to automatically create a static HTML photo gallery from a folder of images.

## Features

- 📸 Automatic detection of image files (JPG, PNG, WebP, GIF)
- 🎨 Responsive masonry layout with CSS Grid
- 🔍 Clickable images for full view
- 📦 Automatically creates a ZIP archive of all photos
- ⚡ Lazy loading for better performance
- 🎯 No external dependencies (only Go and `zip` command)

## Usage

```bash
./gallerygen <image-folder> <output.html>
```

### Example

```bash
./gallerygen ./my-photos gallery.html
```

This creates:
- `gallery.html` - The HTML gallery in the parent directory of the image folder
- `photos.zip` - A ZIP archive of all images in the same directory

### File Structure

```
/home/user/
├── my-photos/            ← Image folder
│   ├── image1.jpg
│   ├── image2.png
│   └── image3.webp
├── gallery.html          ← Generated gallery
└── photos.zip            ← ZIP archive
```

## Supported Image Formats

- JPEG (.jpg, .jpeg)
- PNG (.png)
- WebP (.webp)
- GIF (.gif)

## License

MIT License