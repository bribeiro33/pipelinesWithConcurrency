package imageprocessing

import (
	"image"
	"os"
	"path/filepath"
	"testing"
)

// helpers
func loadTestImage(t *testing.T, path string) image.Image {
	t.Helper()
	img, err := ReadImage(path)
	if err != nil {
		t.Fatalf("failed to read image %q: %v", path, err)
	}
	return img
}

func dimensions(img image.Image) (int, int) {
	b := img.Bounds()
	return b.Dx(), b.Dy()
}

// aspect ratio & resizing testing
func TestResize_Panoramic(t *testing.T) {
	path := "../images/panoramic.jpeg"
	if _, err := os.Stat(path); err != nil {
		t.Skipf("Skipping: missing %q", path)
	}
	img := loadTestImage(t, path)
	resized := Resize(img)
	w, h := dimensions(resized)
	if w > 500 || h > 500 {
		t.Errorf("resized dimensions exceed limit: got %dx%d", w, h)
	}
	if absRatioDiff(img, resized) > 0.02 {
		t.Errorf("aspect ratio changed too much (orig %.3f vs new %.3f)", ratio(img), ratio(resized))
	}
}

func TestResize_Tall(t *testing.T) {
	path := "../images/tall.jpg"
	if _, err := os.Stat(path); err != nil {
		t.Skipf("Skipping: missing %q", path)
	}
	img := loadTestImage(t, path)
	resized := Resize(img)
	w, h := dimensions(resized)
	if w > 500 || h > 500 {
		t.Errorf("resized dimensions exceed limit: got %dx%d", w, h)
	}
	if absRatioDiff(img, resized) > 0.02 {
		t.Errorf("aspect ratio changed too much (orig %.3f vs new %.3f)", ratio(img), ratio(resized))
	}
}

func TestResize_Small_NoChange(t *testing.T) {
	path := "../images/small.png" // 200x200 example
	if _, err := os.Stat(path); err != nil {
		t.Skipf("Skipping: missing %q", path)
	}
	img := loadTestImage(t, path)
	resized := Resize(img)
	origW, origH := dimensions(img)
	newW, newH := dimensions(resized)
	if origW != newW || origH != newH {
		t.Errorf("expected no resize for small image: orig=%dx%d new=%dx%d", origW, origH, newW, newH)
	}
}

// helpers ratio comparison
func ratio(img image.Image) float64 {
	b := img.Bounds()
	return float64(b.Dx()) / float64(b.Dy())
}

func absRatioDiff(a, b image.Image) float64 {
	diff := ratio(a) - ratio(b)
	if diff < 0 {
		diff = -diff
	}
	return diff
}

// transparency and grayscale testing
func TestGrayscale_PNGAlpha(t *testing.T) {
	path := "../images/transparent.png"
	if _, err := os.Stat(path); err != nil {
		t.Skipf("Skipping: missing %q", path)
	}
	img := loadTestImage(t, path)
	gray := Grayscale(img)
	_, ok := gray.(*image.Gray)
	if !ok {
		t.Errorf("expected gray image after Grayscale(), got %T", gray)
	}
}

// error handling for bad inputs testing
func TestReadImage_InvalidFile(t *testing.T) {
	path := "../images/not_an_image.txt"
	if _, err := os.Stat(path); err != nil {
		t.Skipf("Skipping: missing %q", path)
	}
	_, err := ReadImage(path)
	if err == nil {
		t.Errorf("expected error for invalid file type, got nil")
	}
}

func TestReadImage_NonExistent(t *testing.T) {
	path := "../images/does_not_exist.jpg"
	_, err := ReadImage(path)
	if err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}
}

// path handling and filenames testing
func TestWriteImage_SpecialCharacters(t *testing.T) {
	path := "../images/output/🐈cat.jpg"
	os.MkdirAll(filepath.Dir(path), 0o755)
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	if err := WriteImage(path, img); err != nil {
		t.Errorf("WriteImage failed for unicode/space path: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist at %q: %v", path, err)
	}
	// cleanup
	os.Remove(path)
}

// round trip testing
func TestJPEGReadWriteRoundTrip(t *testing.T) {
	path := "../images/cat1.jpg"
	if _, err := os.Stat(path); err != nil {
		t.Skipf("Skipping: missing %q", path)
	}
	img := loadTestImage(t, path)
	out := "../images/output/cat1_test.jpg"
	os.MkdirAll(filepath.Dir(out), 0o755)
	if err := WriteImage(out, img); err != nil {
		t.Fatalf("WriteImage failed: %v", err)
	}
	outImg := loadTestImage(t, out)
	if outImg.Bounds() != img.Bounds() {
		t.Errorf("bounds mismatch after round-trip: orig=%v new=%v", img.Bounds(), outImg.Bounds())
	}

	os.Remove(out)
}
