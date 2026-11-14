// Package media memproses video & foto produk dengan ffmpeg:
//   - video → MP4 H.264, sisi terpanjang ≤ 1280 px, CRF 28, audio AAC 96k, "faststart" (bisa diputar sebelum selesai diunduh)
//   - gambar sampul (poster) dari video, ≤ 720 px
//   - foto → JPEG ≤ 1600 px, orientasi HP diperbaiki, EXIF (termasuk lokasi GPS) dibuang
//
// Butuh ffmpeg + ffprobe di PATH (Mac: brew install ffmpeg). Tanpa itu, file asli dipakai apa adanya.
package media

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	MaxVideoSeconds = 180
	videoMaxSide    = 1280
	posterMaxSide   = 720
	photoMaxSide    = 1600
)

var (
	ErrTooLong = errors.New("video lebih dari 3 menit")
	ErrCorrupt = errors.New("file media rusak")
)

// Available: ffmpeg & ffprobe terpasang.
func Available() bool {
	_, e1 := exec.LookPath("ffmpeg")
	_, e2 := exec.LookPath("ffprobe")
	return e1 == nil && e2 == nil
}

type Info struct {
	Duration float64
	Width    int
	Height   int
	HasAudio bool
	HasVideo bool
}

// Probe membaca durasi, ukuran (sudah memperhitungkan rotasi HP), dan ada/tidaknya audio.
func Probe(ctx context.Context, path string) (Info, error) {
	out, err := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-print_format", "json", "-show_format", "-show_streams", path).Output()
	if err != nil {
		return Info{}, fmt.Errorf("ffprobe: file tidak bisa dibaca sebagai video/gambar")
	}
	var p struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
		Streams []struct {
			CodecType    string            `json:"codec_type"`
			Width        int               `json:"width"`
			Height       int               `json:"height"`
			Tags         map[string]string `json:"tags"`
			SideDataList []struct {
				Rotation float64 `json:"rotation"`
			} `json:"side_data_list"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out, &p); err != nil {
		return Info{}, err
	}
	var in Info
	in.Duration, _ = strconv.ParseFloat(p.Format.Duration, 64)
	for _, s := range p.Streams {
		switch s.CodecType {
		case "audio":
			in.HasAudio = true
		case "video":
			if in.HasVideo {
				continue
			}
			in.HasVideo = true
			in.Width, in.Height = s.Width, s.Height
			rot := 0.0
			if r, err := strconv.ParseFloat(s.Tags["rotate"], 64); err == nil {
				rot = r
			}
			for _, sd := range s.SideDataList {
				if sd.Rotation != 0 {
					rot = sd.Rotation
				}
			}
			if r := int(rot) % 180; r == 90 || r == -90 {
				in.Width, in.Height = in.Height, in.Width
			}
		}
	}
	if !in.HasVideo {
		return in, fmt.Errorf("tidak ada gambar/video di file ini")
	}
	return in, nil
}

// scaleFilter: perkecil supaya sisi terpanjang ≤ max (tidak memperbesar), ukuran genap untuk H.264.
func scaleFilter(max int) string {
	m := strconv.Itoa(max)
	return "scale=w='if(gte(iw,ih),min(" + m + ",iw),-2)':h='if(gte(iw,ih),-2,min(" + m + ",ih))'"
}

func run(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", append([]string{"-hide_banner", "-loglevel", "error", "-y"}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if len(msg) > 400 {
			msg = msg[len(msg)-400:]
		}
		// Waktu habis / dibatalkan → boleh dicoba lagi. Selain itu ffmpeg gagal karena isi file
		// (rusak / format aneh) — mengulang tidak akan membantu. Tidak bergantung pada teks pesan
		// karena berbeda antar versi ffmpeg.
		if ctx.Err() != nil {
			return fmt.Errorf("ffmpeg: %w", ctx.Err())
		}
		return fmt.Errorf("%w: ffmpeg tidak bisa memproses file ini: %v: %s", ErrCorrupt, err, msg)
	}
	return nil
}

// TranscodeVideo: MP4 kecil yang cepat diputar di HP (rotasi HP diterapkan otomatis oleh ffmpeg).
func TranscodeVideo(ctx context.Context, in, out string) error {
	return run(ctx, "-i", in, "-t", strconv.Itoa(MaxVideoSeconds+1),
		"-map", "0:v:0", "-map", "0:a:0?",
		"-vf", scaleFilter(videoMaxSide), "-fpsmax", "30",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "28", "-profile:v", "high", "-pix_fmt", "yuv420p",
		"-c:a", "aac", "-b:a", "96k", "-ac", "2",
		"-movflags", "+faststart", "-map_metadata", "-1", out)
}

// Poster: satu frame (detik ke-1, atau sepertiga durasi untuk video sangat pendek).
func Poster(ctx context.Context, in, out string, duration float64) error {
	at := 1.0
	if duration > 0 && duration < 3 {
		at = duration / 3
	}
	return run(ctx, "-ss", strconv.FormatFloat(at, 'f', 2, 64), "-i", in, "-frames:v", "1",
		"-vf", scaleFilter(posterMaxSide), "-q:v", "3", out)
}

// ResizePhoto: JPEG ≤ 1600 px; orientasi dari EXIF diterapkan, metadata (termasuk GPS) dibuang.
func ResizePhoto(ctx context.Context, in, out string) error {
	vf := scaleFilter(photoMaxSide)
	if t := orientationFilter(jpegOrientation(in)); t != "" {
		vf = t + "," + vf
	}
	return run(ctx, "-i", in, "-frames:v", "1", "-vf", vf, "-q:v", "3", "-map_metadata", "-1", out)
}

func orientationFilter(o int) string {
	return map[int]string{2: "hflip", 3: "hflip,vflip", 4: "vflip", 5: "transpose=0", 6: "transpose=1", 7: "transpose=3", 8: "transpose=2"}[o]
}

// jpegOrientation membaca tag EXIF Orientation (0x0112) dari JPEG. 1/0 = normal.
func jpegOrientation(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	buf := make([]byte, 128<<10)
	n, _ := io.ReadFull(f, buf)
	return exifOrientation(buf[:n])
}

func exifOrientation(b []byte) int {
	if len(b) < 4 || b[0] != 0xFF || b[1] != 0xD8 {
		return 0
	}
	i := 2
	for i+4 <= len(b) {
		if b[i] != 0xFF {
			return 0
		}
		marker := b[i+1]
		size := int(binary.BigEndian.Uint16(b[i+2:]))
		if marker == 0xE1 && i+4+size <= len(b) && size > 14 && string(b[i+4:i+10]) == "Exif\x00\x00" {
			t := b[i+10 : i+2+size]
			var bo binary.ByteOrder
			switch string(t[:2]) {
			case "II":
				bo = binary.LittleEndian
			case "MM":
				bo = binary.BigEndian
			default:
				return 0
			}
			ifd := int(bo.Uint32(t[4:]))
			if ifd+2 > len(t) {
				return 0
			}
			count := int(bo.Uint16(t[ifd:]))
			for k := 0; k < count; k++ {
				e := ifd + 2 + k*12
				if e+12 > len(t) {
					return 0
				}
				if bo.Uint16(t[e:]) == 0x0112 {
					return int(bo.Uint16(t[e+8:]))
				}
			}
			return 0
		}
		if marker == 0xDA { // mulai data gambar: tidak ada EXIF
			return 0
		}
		i += 2 + size
	}
	return 0
}
