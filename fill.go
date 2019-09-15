package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/kettek/apng"
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

func fillAPNGBackground(si StickerImage, clr color.Color) StickerImage {
	imgs := si.raw.(apng.APNG)
	frames := []apng.Frame{}
	for _, frame := range imgs.Frames {
		b := frame.Image.Bounds()

		out := image.NewRGBA(b)
		for x := b.Min.X; x < b.Max.X; x++ {
			for y := b.Min.Y; y < b.Max.Y; y++ {
				out.Set(x, y, clr)
			}
		}

		draw.Draw(out, b, frame.Image, image.Pt(0, 0), draw.Over)
		frames = append(frames, apng.Frame{Image: out})
	}

	return StickerImage{
		Type: si.Type,
		raw: apng.APNG{
			Frames:    frames,
			LoopCount: 0,
		},
	}
}
