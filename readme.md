# Gallery Generator

A robust Go tool to automatically create a static HTML photo gallery from a folder of images.

## Features

- 📸 **Automatic Detection**: Finds JPG, PNG, WebP, and GIF files.
- 🖼️ **Thumbnails**: Automatically generates optimized thumbnails for fast loading.
- 🧱 **Masonry Layout**: Beautiful, responsive layout that handles images of varying heights (like Pinterest).
- 🔍 **Lightbox**: Clickable images open in a full-screen overlay with keyboard navigation (Left/Right arrows).
- 📦 **Zip Archive**: Automatically creates a downloadable ZIP of all photos.
- ⚡ **Performance**: Lazy loading and thumbnail usage ensure the gallery loads instantly.
- 🛠️ **Zero Dependencies**: Native Go implementation (no external `zip` command required).
- 🖥️ **Interactive Mode**: Simply double-click to run, or use CLI arguments.

## Usage

### Interactive Mode (Recommended)

Just run the executable without arguments:

```bash
./gallerygen
```

It will prompt you for:
1. Image folder path
2. Gallery title
3. Output filename
4. Whether to create a ZIP archive

### CLI Mode

You can also use command-line flags for automation:

```bash
./gallerygen [flags] <image-folder>
```

#### Flags

- `-title string`: Title of the gallery (default "Photos")
- `-output string`: Output HTML filename (default "gallery.html")
- `-no-zip`: Skip creating the zip archive

#### Example

```bash
./gallerygen -title "Holiday 2024" -output index.html ./my-photos
```

This creates:
- `index.html` - The HTML gallery
- `thumbnails/` - Folder containing generated thumbnails
- `photos.zip` - A ZIP archive of all images (unless `-no-zip` is used)

### File Structure

```
/home/user/
├── my-photos/            ← Image folder
│   ├── image1.jpg
│   ├── image2.png
│   └── image3.webp
├── thumbnails/           ← Generated thumbnails
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