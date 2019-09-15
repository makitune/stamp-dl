package main

import (
	"image"
	"image/color"
	"image/draw"
)

func fillImageBackground(si StickerImage, clr color.Color) StickerImage {
	img := si.raw.(image.Image)
	b := img.Bounds()

	out := image.NewRGBA(b)
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			out.Set(x, y, clr)
		}
	}

	draw.Draw(out, b, img, image.Pt(0, 0), draw.Over)
	return StickerImage{
		Type: si.Type,
		raw:  out,
	}
}
