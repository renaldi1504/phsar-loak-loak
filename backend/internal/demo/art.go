package demo

// Gambar contoh (tanpa file dari luar): latar gradasi + bentuk sederhana sesuai jenis barang.
// Hanya untuk data demo — seller asli mengunggah foto & video produknya sendiri.

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"math"
)

type rgb struct{ r, g, b float64 }

func hex(c uint32) rgb {
	return rgb{float64(c >> 16 & 0xff), float64(c >> 8 & 0xff), float64(c & 0xff)}
}

func mix(a, b rgb, t float64) rgb {
	return rgb{a.r + (b.r-a.r)*t, a.g + (b.g-a.g)*t, a.b + (b.b-a.b)*t}
}

var white = rgb{255, 255, 255}

type canvas struct {
	w, h int
	px   []rgb
}

func newCanvas(w, h int) *canvas { return &canvas{w, h, make([]rgb, w*h)} }

func (c *canvas) bg(top, bottom rgb) {
	for y := 0; y < c.h; y++ {
		col := mix(top, bottom, float64(y)/float64(c.h))
		for x := 0; x < c.w; x++ {
			c.px[y*c.w+x] = col
		}
	}
}

// fill: sdf(x,y) < 0 di dalam bentuk; tepi dihaluskan 1 px.
func (c *canvas) fill(sdf func(x, y float64) float64, col func(x, y float64) rgb, alpha float64) {
	for y := 0; y < c.h; y++ {
		for x := 0; x < c.w; x++ {
			fx, fy := float64(x)+.5, float64(y)+.5
			d := sdf(fx, fy)
			if d > 1 {
				continue
			}
			a := math.Min(1, math.Max(0, .5-d)) * alpha
			i := y*c.w + x
			c.px[i] = mix(c.px[i], col(fx, fy), a)
		}
	}
}

func solid(c rgb) func(x, y float64) rgb { return func(_, _ float64) rgb { return c } }

func vgrad(a, b rgb, y0, y1 float64) func(x, y float64) rgb {
	return func(_, y float64) rgb { return mix(a, b, math.Min(1, math.Max(0, (y-y0)/(y1-y0)))) }
}

func circle(cx, cy, r float64) func(x, y float64) float64 {
	return func(x, y float64) float64 { return math.Hypot(x-cx, y-cy) - r }
}

func ellipse(cx, cy, rx, ry float64) func(x, y float64) float64 {
	return func(x, y float64) float64 {
		k := math.Hypot((x-cx)/rx, (y-cy)/ry)
		return (k - 1) * math.Min(rx, ry)
	}
}

func rrect(cx, cy, hw, hh, r float64) func(x, y float64) float64 {
	return func(x, y float64) float64 {
		qx, qy := math.Abs(x-cx)-hw+r, math.Abs(y-cy)-hh+r
		return math.Hypot(math.Max(qx, 0), math.Max(qy, 0)) + math.Min(math.Max(qx, qy), 0) - r
	}
}

// rotated rounded rect (sudut dalam radian).
func rrectRot(cx, cy, hw, hh, r, ang float64) func(x, y float64) float64 {
	base := rrect(0, 0, hw, hh, r)
	s, co := math.Sin(-ang), math.Cos(-ang)
	return func(x, y float64) float64 {
		dx, dy := x-cx, y-cy
		return base(dx*co-dy*s, dx*s+dy*co)
	}
}

func union(fs ...func(x, y float64) float64) func(x, y float64) float64 {
	return func(x, y float64) float64 {
		m := math.Inf(1)
		for _, f := range fs {
			m = math.Min(m, f(x, y))
		}
		return m
	}
}

func ring(cx, cy, r, w float64) func(x, y float64) float64 {
	return func(x, y float64) float64 { return math.Abs(math.Hypot(x-cx, y-cy)-r) - w/2 }
}

func clipAbove(f func(x, y float64) float64, yMax float64) func(x, y float64) float64 {
	return func(x, y float64) float64 { return math.Max(f(x, y), y-yMax) }
}

func (c *canvas) jpeg(q int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, c.w, c.h))
	for i, p := range c.px {
		img.Pix[i*4], img.Pix[i*4+1], img.Pix[i*4+2], img.Pix[i*4+3] = uint8(p.r), uint8(p.g), uint8(p.b), 255
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: q})
	return buf.Bytes()
}

// productPhoto: view 0 = foto utama, 1 = foto kedua (lebih kecil, sudut lain), 2 = foto kondisi (close-up + goresan).
func productPhoto(shape string, colors [2]uint32, accent uint32, view int) []byte {
	const S = 1000
	c := newCanvas(S, S)
	a, b := hex(colors[0]), hex(colors[1])
	acc := hex(accent)
	c.bg(mix(acc, white, .86), mix(acc, white, .7))
	scale := 1.0
	cx, cy := 500.0, 520.0
	if view == 1 {
		scale, cx, cy = .78, 560, 540
		c.bg(mix(acc, white, .92), mix(b, white, .75))
	}
	if view == 2 {
		scale, cx, cy = 1.7, 470, 560
	}
	s := func(v float64) float64 { return v * scale }
	// bayangan lembut
	c.fill(ellipse(cx, cy+s(330), s(300), s(40)), solid(rgb{0, 0, 0}), .12)
	body := vgrad(mix(a, white, .15), b, cy-s(350), cy+s(350))
	shine := solid(white)
	switch shape {
	case "phone":
		c.fill(rrect(cx, cy, s(190), s(360), s(48)), solid(mix(a, rgb{0, 0, 0}, .35)), 1)
		c.fill(rrect(cx, cy, s(172), s(342), s(36)), body, 1)
		c.fill(rrect(cx, cy-s(318), s(46), s(10), s(10)), solid(mix(a, rgb{0, 0, 0}, .5)), 1)
		c.fill(rrectRot(cx-s(60), cy-s(80), s(30), s(330), s(20), -.5), shine, .12)
	case "round":
		c.fill(circle(cx, cy, s(270)), body, 1)
		c.fill(ring(cx, cy, s(200), s(26)), solid(mix(b, white, .35)), .8)
		c.fill(circle(cx, cy, s(90)), solid(mix(a, rgb{0, 0, 0}, .3)), 1)
		c.fill(circle(cx-s(110), cy-s(120), s(60)), shine, .25)
	case "bottle":
		c.fill(rrect(cx, cy-s(330), s(70), s(50), s(14)), solid(mix(b, rgb{0, 0, 0}, .25)), 1)
		c.fill(rrect(cx, cy+s(40), s(170), s(300), s(60)), body, 1)
		c.fill(rrect(cx, cy+s(60), s(150), s(110), s(16)), solid(mix(white, a, .1)), .92)
		c.fill(rrect(cx, cy+s(40), s(90), s(14), s(7)), solid(b), .9)
		c.fill(rrect(cx, cy+s(85), s(60), s(10), s(5)), solid(mix(b, white, .4)), .9)
		c.fill(rrectRot(cx-s(110), cy-s(40), s(16), s(180), s(16), 0), shine, .25)
	case "shirt":
		shirt := union(rrect(cx, cy+s(60), s(200), s(270), s(30)),
			rrectRot(cx-s(240), cy-s(130), s(70), s(120), s(26), .9), rrectRot(cx+s(240), cy-s(130), s(70), s(120), s(26), -.9),
			rrect(cx, cy-s(200), s(200), s(40), s(30)))
		c.fill(shirt, body, 1)
		c.fill(ellipse(cx, cy-s(238), s(70), s(34)), solid(mix(acc, white, .8)), 1)
		c.fill(rrect(cx, cy+s(10), s(4), s(240), s(2)), solid(mix(b, rgb{0, 0, 0}, .2)), .3)
	case "pan":
		c.fill(rrectRot(cx+s(300), cy+s(170), s(160), s(26), s(24), .55), solid(mix(a, rgb{0, 0, 0}, .2)), 1)
		c.fill(circle(cx-s(40), cy-s(20), s(270)), solid(mix(a, rgb{0, 0, 0}, .3)), 1)
		c.fill(circle(cx-s(40), cy-s(20), s(235)), body, 1)
		c.fill(circle(cx-s(110), cy-s(90), s(70)), shine, .14)
	case "bag":
		c.fill(clipAbove(ring(cx, cy-s(120), s(150), s(30)), cy-s(120)), solid(mix(b, rgb{0, 0, 0}, .25)), 1)
		c.fill(rrect(cx, cy+s(110), s(270), s(220), s(50)), body, 1)
		c.fill(rrect(cx, cy+s(20), s(270), s(12), s(6)), solid(mix(b, rgb{0, 0, 0}, .2)), .6)
		c.fill(rrect(cx, cy+s(60), s(40), s(26), s(8)), solid(rgb{212, 175, 55}), 1)
	case "cap":
		c.fill(ellipse(cx+s(120), cy+s(110), s(300), s(70)), solid(mix(b, rgb{0, 0, 0}, .15)), 1)
		c.fill(clipAbove(circle(cx, cy+s(110), s(250)), cy+s(110)), body, 1)
		c.fill(circle(cx, cy-s(135), s(22)), solid(mix(a, rgb{0, 0, 0}, .3)), 1)
		c.fill(rrect(cx, cy-s(10), s(80), s(40), s(10)), solid(mix(white, b, .1)), .85)
	default: // box
		c.fill(rrect(cx, cy, s(270), s(270), s(40)), body, 1)
		c.fill(rrectRot(cx, cy, s(400), s(36), s(10), -.78), solid(mix(b, white, .35)), .55)
		c.fill(rrect(cx, cy+s(170), s(140), s(26), s(10)), solid(mix(white, a, .1)), .85)
	}
	if view == 2 { // goresan untuk foto kondisi barang bekas
		for i := 0; i < 3; i++ {
			c.fill(rrectRot(420+float64(i)*70, 430+float64(i)*40, 120-float64(i)*25, 2.2, 1, -.35+float64(i)*.12), solid(white), .7)
		}
		c.fill(ring(420, 430, 190, 8), solid(rgb{220, 38, 38}), .9)
	}
	return c.jpeg(86)
}

func storeLogo(colors [2]uint32) []byte {
	c := newCanvas(400, 400)
	a, b := hex(colors[0]), hex(colors[1])
	c.bg(a, b)
	c.fill(circle(200, 200, 120), solid(white), .95)
	c.fill(circle(200, 200, 72), vgrad(a, b, 130, 270), 1)
	c.fill(circle(170, 168, 22), solid(white), .5)
	return c.jpeg(88)
}

func storeBanner(colors [2]uint32) []byte {
	c := newCanvas(1600, 500)
	a, b := hex(colors[0]), hex(colors[1])
	for y := 0; y < c.h; y++ {
		for x := 0; x < c.w; x++ {
			c.px[y*c.w+x] = mix(a, b, (float64(x)/1600+float64(y)/500)/2)
		}
	}
	for i := 0; i < 9; i++ {
		f := float64(i)
		c.fill(circle(120+f*190, 250+math.Sin(f)*150, 60+math.Mod(f*37, 90)), solid(white), .1)
	}
	return c.jpeg(84)
}

var _ = color.RGBA{}
