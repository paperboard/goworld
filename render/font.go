package render

import (
	"image"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	File    string
	Size    float32
	DPI     float32
	Spacing float32
	Color   Color

	fnt    *truetype.Font
	drawer *font.Drawer
}

func (f *Font) setup() {
	f.drawer = &font.Drawer{
		Face: truetype.NewFace(f.fnt, &truetype.Options{
			Size:    float64(f.Size),
			DPI:     float64(72 * f.DPI),
			Hinting: font.HintingFull,
		}),
	}
}

func (f *Font) LineHeight() float32 {
	return math.Ceil(f.Size * f.Spacing * f.DPI)
}

func (f *Font) Measure(text string) vec2.T {
	lines := 1
	width := 0
	s := 0
	for i, c := range text {
		if c == '\n' {
			line := text[s:i]
			w := f.drawer.MeasureString(line).Ceil()
			if w > width {
				width = w
			}
			s = i + 1
			lines++
		}
	}
	r := len(text)
	if s < r {
		line := text[s:]
		w := f.drawer.MeasureString(line).Ceil()
		if w > width {
			width = w
		}
	}

	lineHeight := int(f.LineHeight())
	height := lineHeight*lines + (lineHeight / 2)
	return vec2.NewI(width, height)
}

func (f *Font) RenderNew(text string, color Color) *Texture {
	size := f.Measure(text)
	texture := CreateTexture(int(size.X), int(size.Y))
	f.Render(texture, text, color)
	return texture
}

func (f *Font) Render(tx *Texture, text string, color Color) {
	f.drawer.Src = image.NewUniform(color.RGBA())

	size := f.Measure(text)

	// todo: its probably not a great idea to allocate an image on every draw
	// perhaps textures should always have a backing image ?
	rgba := image.NewRGBA(image.Rect(0, 0, int(math.Ceil(size.X)), int(math.Ceil(size.Y))))
	f.drawer.Dst = rgba

	s := 0
	line := 1
	lineHeight := int(f.LineHeight())
	for i, c := range text {
		if c == '\n' {
			if i == s {
				continue // skip empty rows
			}
			f.drawer.Dot = fixed.P(0, line*int(lineHeight))
			f.drawer.DrawString(text[s:i])
			s = i + 1
			line++
		}
	}
	if s < len(text) {
		f.drawer.Dot = fixed.P(0, line*int(lineHeight))
		f.drawer.DrawString(text[s:])
	}

	tx.Bind()
	tx.Buffer(rgba)
}

/** Load a truetype font */
func LoadFont(filename string, size, dpi, spacing float32) *Font {
	fontBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	fnt := &Font{
		Size:    size,
		DPI:     dpi,
		Spacing: spacing,
		fnt:     f,
	}
	fnt.setup()
	return fnt
}
