package main

import (
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"os"
	"path/filepath"
	"testing"
)

func TestRunSequential(t *testing.T) {
	inputs := []string{
		"images/large.jpg",
		"images/cat1.jpg",
		"images/cat2.jpg",
	}

	// check input files actually exist
	for _, p := range inputs {
		if _, err := os.Stat(p); err != nil {
			t.Skipf("Skipping test: input file %q not found", p)
		}
	}

	if err := RunSequential(inputs); err != nil {
		t.Fatalf("RunSequential failed: %v", err)
	}

	for _, p := range inputs {
		out := filepath.Join("images", "output", filepath.Base(p))
		if _, err := os.Stat(out); err != nil {
			t.Errorf("Output file %q not found: %v", out, err)
		}
	}
}

func TestRunConcurrent(t *testing.T) {
	inputs := []string{
		"images/large.jpg",
		"images/cat1.jpg",
		"images/cat2.jpg",
	}

	for _, p := range inputs {
		if _, err := os.Stat(p); err != nil {
			t.Skipf("Skipping test: input file %q not found", p)
		}
	}

	if err := RunConcurrent(inputs); err != nil {
		t.Fatalf("RunConcurrent failed: %v", err)
	}

	for _, p := range inputs {
		out := filepath.Join("images", "output", filepath.Base(p))
		if _, err := os.Stat(out); err != nil {
			t.Errorf("Output file %q not found: %v", out, err)
		}
	}
}

// Benchmark the sequential pipeline
func BenchmarkSequential(b *testing.B) {
	images := []string{
		"images/cat1.jpg",
		"images/cat2.jpg",
	}

	// Pre-check that inputs exist, otherwise skip benchmark
	for _, p := range images {
		if _, err := os.Stat(p); err != nil {
			b.Skipf("Skipping benchmark: input file %q not found", p)
		}
	}

	b.ResetTimer() // don’t include setup time in the measurement
	for i := 0; i < b.N; i++ {
		if err := RunSequential(images); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark the concurrent pipeline
func BenchmarkConcurrent(b *testing.B) {
	images := []string{
		"images/cat1.jpg",
		"images/cat2.jpg",
	}

	for _, p := range images {
		if _, err := os.Stat(p); err != nil {
			b.Skipf("Skipping benchmark: input file %q not found", p)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RunConcurrent(images); err != nil {
			b.Fatal(err)
		}
	}
}

func loadTestImage(t *testing.T, path string) image.Image {
	t.Helper()
	img, err := imageprocessing.ReadImage(path)
	if err != nil {
		t.Fatalf("failed to read image %q: %v", path, err)
	}
	return img
}

func BenchmarkStage_Load(b *testing.B) {
	imgs := []string{"images/cat2.jpg", "images/cat2.jpg"}
	for _, p := range imgs {
		if _, err := os.Stat(p); err != nil {
			b.Skip("missing test image")
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range imgs {
			_, _ = imageprocessing.ReadImage(p)
		}
	}
}

func BenchmarkStage_Resize(b *testing.B) {
	img := loadTestImage(&testing.T{}, "images/cat2.jpg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = imageprocessing.Resize(img)
	}
}

func BenchmarkStage_Grayscale(b *testing.B) {
	img := loadTestImage(&testing.T{}, "images/cat2.jpg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = imageprocessing.Grayscale(img)
	}
}

func BenchmarkStage_Write(b *testing.B) {
	img := loadTestImage(&testing.T{}, "images/cat2.jpg")
	out := "images/output/bench_write.jpg"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = imageprocessing.WriteImage(out, img)
	}
}

func BenchmarkConcurrentScaling(b *testing.B) {
	imageSets := [][]string{
		{"images/cat1.jpg"},
		{"images/cat1.jpg", "images/cat2.jpg"},
		{"images/cat1.jpg", "images/cat2.jpg", "images/panoramic.jpeg", "images/tall.jpg"},
	}
	for _, imgs := range imageSets {
		b.Run(fmt.Sprintf("%d_images", len(imgs)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if err := RunConcurrent(imgs); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
