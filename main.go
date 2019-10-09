// http://fabiensanglard.net/doom_fire_psx/

package main

import (
	"errors"
	"flag"
	"image"
	"image/draw"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type drawContext struct {
	img   *image.Paletted
	off   bool
	black int
	white int
}

func newDrawContext(x, y int) *drawContext {
	rand.Seed(time.Now().UnixNano())
	dc := drawContext{
		img:   image.NewPaletted(image.Rect(0, 0, x, y), palette),
		black: 0,
		white: len(palette) - 1,
	}
	seed(dc.img, dc.white)
	return &dc
}

func seed(img *image.Paletted, c int) {
	r := img.Bounds().Max
	for x := 0; x < r.X; x++ {
		img.SetColorIndex(x, r.Y-1, uint8(c))
	}
}

func (dc *drawContext) toggle() {
	color := dc.white
	if dc.off = !dc.off; dc.off {
		color = dc.black
	}
	seed(dc.img, color)
}

func (dc *drawContext) drawTo(dst draw.Image) {
	r := dc.img.Bounds().Max
	for x := 0; x < r.X; x++ {
		for y := r.Y - 1; y > 0; y-- {
			z := rand.Intn(3) - 1 // -1, 0, 1
			n := dc.img.ColorIndexAt(x, y)
			if n > 0 && z == 0 {
				n-- // next color
			}
			dc.img.SetColorIndex(x+z, y-1, n)
		}
	}
	draw.Draw(dst, dst.Bounds(), dc.img, image.ZP, draw.Src)
}

func (dc *drawContext) update(screen *ebiten.Image) error {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyQ):
		return errors.New("exit")
	case inpututil.IsKeyJustPressed(ebiten.KeySpace):
		dc.toggle()
	}
	if !ebiten.IsDrawingSkipped() {
		dc.drawTo(screen)
	}
	return nil
}

func main() {
	width := flag.Int("width", 320, "screen width")
	height := flag.Int("height", 200, "screen height")
	scale := flag.Float64("scale", 2.0, "scale")
	flag.Parse()

	dc := newDrawContext(*width, *height)
	ebiten.SetRunnableInBackground(true)
	if err := ebiten.Run(dc.update, *width, *height, *scale, "Fire"); err != nil {
		log.Fatal(err)
	}
}
