package main

import (
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"strings"
)

// sequential version (no goroutines)
func RunSequential(imagePaths []string) error {
	for _, p := range imagePaths {
		out := strings.Replace(p, "images/", "images/output/", 1)
		img, err := imageprocessing.ReadImage(p)
		if err != nil {
			return err
		}
		img = imageprocessing.Resize(img)
		img = imageprocessing.Grayscale(img)
		if err := ensureOutDir(out); err != nil {
			return err
		}
		if err := imageprocessing.WriteImage(out, img); err != nil {
			return err
		}
	}
	return nil
}

// concurrent/pipeline version
func RunConcurrent(imagePaths []string) error {
	c1 := loadImage(imagePaths)
	c2 := resize(c1)
	c3 := convertToGrayscale(c2)
	results := saveImage(c3)

	var hadErr bool
	for err := range results {
		if err != nil {
			hadErr = true
		}
	}
	if hadErr {
		return fmt.Errorf("some images failed to process")
	}
	return nil
}
