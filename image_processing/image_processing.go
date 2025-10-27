package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// reads an image file from disk and returns the decoded image
func ReadImage(path string) (image.Image, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image %q: %w", path, err)
	}
	defer inputFile.Close()

	// Decode the image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, fmt.Errorf("decode image %q: %w", path, err)
	}
	return img, nil
}

// writes an image to disk in JPEG
func WriteImage(path string, img image.Image) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create image %q: %w", path, err)
	}
	defer outputFile.Close()

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		opts := &jpeg.Options{Quality: 90}
		if err := jpeg.Encode(outputFile, img, opts); err != nil {
			return fmt.Errorf("jpeg encode %q: %w", path, err)
		}
	case ".png":
		if err := png.Encode(outputFile, img); err != nil {
			return fmt.Errorf("png encode %q: %w", path, err)
		}
	default:
		return fmt.Errorf("unsupported image format: %s (only .jpg/.jpeg/.png supported)", ext)
	}

	return nil
}

// converts image to greyscale
func Grayscale(img image.Image) image.Image {
	// Create a new grayscale image
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Convert each pixel to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			grayPixel := color.GrayModel.Convert(originalPixel)
			grayImg.Set(x, y, grayPixel)
		}
	}
	return grayImg
}

// resizes image to max 500,500 while retaining original dimensions
func Resize(img image.Image) image.Image {
	const maxWidth, maxHeight = 500, 500

	// original bounds
	bounds := img.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	// determine scale factor that fits within the bounding box
	scale := math.Min(float64(maxWidth)/width, float64(maxHeight)/height)

	// only resize if the image is larger than the max dimensions
	if scale >= 1.0 {
		return img // already small
	}

	newWidth := uint(width * scale)
	newHeight := uint(height * scale)

	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}
