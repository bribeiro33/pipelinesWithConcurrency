package main

import (
	"flag"
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"os"
	"path/filepath"
	"strings"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
	Err       error
}

func ensureOutDir(p string) error {
	return os.MkdirAll(filepath.Dir(p), 0o755)
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		defer close(out)
		// For each input path create a job and add it to
		// the out channel
		for _, p := range paths {
			job := Job{InputPath: p,
				OutPath: strings.Replace(p, "images/", "images/output/", 1)}
			img, err := imageprocessing.ReadImage(p)
			if err != nil {
				job.Err = err
			} else {
				job.Image = img
			}
			out <- job
		}
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		defer close(out)
		// For each input job, create a new job after resize and add it to
		// the out channel
		for job := range input { // Read from the channel
			if job.Err == nil {
				job.Image = imageprocessing.Resize(job.Image)
			}
			out <- job
		}
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		defer close(out)
		for job := range input { // Read from the channel
			if job.Err == nil {
				job.Image = imageprocessing.Grayscale(job.Image)
			}
			out <- job
		}
	}()
	return out
}

func saveImage(input <-chan Job) <-chan error {
	out := make(chan error)
	go func() {
		defer close(out)
		for job := range input { // Read from the channel
			if job.Err != nil {
				out <- job.Err
				continue
			}
			if err := ensureOutDir(job.OutPath); err != nil {
				out <- err
				continue
			}
			if err := imageprocessing.WriteImage(job.OutPath, job.Image); err != nil {
				out <- err
				continue
			}
			out <- nil
		}
	}()
	return out
}

func main() {
	// CLI input --images="img1.jpg, img2.png,"
	pathsFlag := flag.String("images", "", "Comma-separated list of image paths to process")
	flag.Parse()

	if *pathsFlag == "" {
		fmt.Fprintln(os.Stderr, "Error: no image paths provided.\nUsage: go run . --images=\"path1.jpg,path2.png\"")
		os.Exit(1)
	}

	imagePaths := strings.Split(*pathsFlag, ",")
	fmt.Println("Processing images:", imagePaths)

	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	results := saveImage(channel3)

	var failed bool
	for err := range results {
		if err != nil {
			failed = true
			fmt.Println("Failed: ", err)
		} else {
			fmt.Println("Success!")
		}
	}
	if failed {
		os.Exit(1)
	}

}
