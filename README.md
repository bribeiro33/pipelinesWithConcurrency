# Go Image Processing Pipeline

This project builds upon amritsingh's project (https://github.com/code-heim/go_21_goroutines_pipeline) and implements a concurrent image-processing pipeline in Go that loads images, resizes them while preserving aspect ratio, converts them to grayscale, and writes the results to an output folder with robust error handling.

It demonstrates:

* Go goroutines and channels for pipelined concurrency

* Proper error handling and I/O safety

* Unit testing and benchmarking (stage-level and end-to-end)

* Aspect-ratio-preserving resizing and format-aware image writing (JPEG/PNG)

Requires Golang 1.20 or higher

## Structure
```console
goroutines_pipeline/
├── main.go                        # Defines Job struct and pipeline stages (load, resize, grayscale, save)
├── main_test.go                   # End-to-end tests and benchmarks (sequential vs concurrent)
├── image_processing/
│   ├── image_processing.go        # Main image helpers: ReadImage, WriteImage, Resize, Grayscale
│   └── image_processing_test.go   # Unit tests for each image-processing function
└── images/                        # Tester images
    ├── cat1.jpg
    ├── cat2.jpg
    ├── transparent.png
    ├── panoramic.jpg
    ├── tall.png
    ├── small.jpg
    └── output/                   # Image output location               

```

## Features
### Image-processing

* ReadImage(path) — reads .jpg or .png (added) files with error checking (added)

* WriteImage(path, img) — writes output in the correct format based on file extension

* Resize(img) — resizes images to fit within 500×500 pixels (added) while maintaining aspect ratio

* Grayscale(img) — converts any image to grayscale safely, including PNGs with transparency

### Pipeline (main.go)

Implements four connected stages using Go channels:

1. Load — read images from disk

2. Resize — scale proportionally

3. Grayscale — convert to black and white

4. Save — write processed images to images/output/

Two pipeline modes are supported:

* RunSequential(images []string) — processes images in order, one at a time

* RunConcurrent(images []string) — runs each stage concurrently using goroutines and channels

### Error handling

* Handles missing or invalid files (no panics).

* Reports unsupported formats (only .jpg and .png allowed).

* Ensures output directories exist before writing.

## How to use
1. Clone repo
2. Run program with images flag e.g.:
```console
go run . --images="images/<img_name>.png,images/<img_name>.jpg"
```
**Output images will be under the original image name in /images/output/**

4. Run all tests 
```console
go test ./... -v
```
4. Run all benchmarks
```console
go test -bench=. -benchmem
```
