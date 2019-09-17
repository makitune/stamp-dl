package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/kettek/apng"
)

func (si *LineStickerImage) fillPNGBackground(clr color.Color) LineStickerImage {
	img := si.raw.(image.Image)
	out := filledImage(img, clr)
	return LineStickerImage{
		Type: si.Type,
		raw:  out,
	}
}

func (si *LineStickerImage) fillAPNGBackground(clr color.Color) LineStickerImage {
	imgs := si.raw.(apng.APNG)
	frames := []apng.Frame{}
	for _, frame := range imgs.Frames {
		out := filledImage(frame.Image, clr)
		img := frame
		img.Image = out
		frames = append(frames, img)
	}

	return LineStickerImage{
		Type: si.Type,
		raw: apng.APNG{
			Frames:    frames,
			LoopCount: 0,
		},
	}
}

func filledImage(img image.Image, clr color.Color) image.Image {
	b := img.Bounds()

	out := image.NewRGBA(b)
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			out.Set(x, y, clr)
		}
	}

	draw.Draw(out, b, img, image.Pt(0, 0), draw.Over)
	return out
}
