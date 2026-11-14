package media

import (
	"context"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

func need(t *testing.T) {
	if !Available() {
		t.Skip("ffmpeg tidak terpasang")
	}
}

// makeVideo membuat video uji: potret 1080x1920, 30 fps, dengan audio.
func makeVideo(t *testing.T, dir string, seconds string) string {
	p := filepath.Join(dir, "in.mp4")
	out, err := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-y",
		"-f", "lavfi", "-i", "testsrc2=size=1080x1920:rate=30", "-f", "lavfi", "-i", "sine=frequency=440",
		"-t", seconds, "-c:v", "libx264", "-preset", "ultrafast", "-c:a", "aac", "-shortest", p).CombinedOutput()
	if err != nil {
		t.Fatalf("buat video uji: %v %s", err, out)
	}
	return p
}

func TestTranscodeAndPoster(t *testing.T) {
	need(t)
	dir := t.TempDir()
	in := makeVideo(t, dir, "3")
	ctx := context.Background()
	out, poster := filepath.Join(dir, "out.mp4"), filepath.Join(dir, "p.jpg")
	if err := TranscodeVideo(ctx, in, out); err != nil {
		t.Fatal(err)
	}
	info, err := Probe(ctx, out)
	if err != nil {
		t.Fatal(err)
	}
	if info.Height != 1280 || info.Width != 720 || !info.HasAudio || info.Duration < 2.5 {
		t.Fatalf("hasil kompres: %+v (harusnya 720x1280 + audio)", info)
	}
	if err := Poster(ctx, out, poster, info.Duration); err != nil {
		t.Fatal(err)
	}
	pi, err := Probe(ctx, poster)
	if err != nil || pi.Height != 720 {
		t.Fatalf("poster: %+v %v", pi, err)
	}
}

func TestResizePhotoKeepsSmallAndShrinksLarge(t *testing.T) {
	need(t)
	dir := t.TempDir()
	ctx := context.Background()
	for _, c := range []struct{ size, want string }{{"4000x3000", "1600x1200"}, {"800x600", "800x600"}} {
		in := filepath.Join(dir, "in.jpg")
		if out, err := exec.Command("ffmpeg", "-loglevel", "error", "-y", "-f", "lavfi", "-i", "testsrc2=size="+c.size, "-frames:v", "1", in).CombinedOutput(); err != nil {
			t.Fatalf("%v %s", err, out)
		}
		out := filepath.Join(dir, "out.jpg")
		if err := ResizePhoto(ctx, in, out); err != nil {
			t.Fatal(err)
		}
		info, _ := Probe(ctx, out)
		if got := itoa(info.Width) + "x" + itoa(info.Height); got != c.want {
			t.Errorf("%s → %s, harusnya %s", c.size, got, c.want)
		}
	}
}

func TestExifOrientation(t *testing.T) {
	// JPEG minimal: SOI + APP1 Exif (little endian) dengan Orientation = 6.
	tiff := []byte{'I', 'I', 42, 0, 8, 0, 0, 0, 1, 0, 0x12, 0x01, 3, 0, 1, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0}
	app1 := append([]byte("Exif\x00\x00"), tiff...)
	n := len(app1) + 2
	b := append([]byte{0xFF, 0xD8, 0xFF, 0xE1, byte(n >> 8), byte(n)}, app1...)
	b = append(b, 0xFF, 0xDA)
	if o := exifOrientation(b); o != 6 {
		t.Fatalf("orientation = %d, harusnya 6", o)
	}
	if orientationFilter(6) != "transpose=1" || orientationFilter(1) != "" {
		t.Fatal("filter rotasi salah")
	}
	if exifOrientation([]byte("bukan jpeg")) != 0 {
		t.Fatal("bukan jpeg harus 0")
	}
}

func itoa(i int) string { return strconv.Itoa(i) }
