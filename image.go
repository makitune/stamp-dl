package main

import (
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/kettek/apng"
)

// PNG is an object for image formated png.
type PNG struct {
	name  string
	image image.Image
}

// Encode fills the background and writes the PNG.image to w in PNG format.
func (p *PNG) Encode(w io.Writer) error {
	img := filledImage(p.image, color.White)
	return png.Encode(w, img)
}

// StoreName return the file name for PNG.
func (p *PNG) StoreName() string {
	return p.name + ".png"
}

// APNG is an object for image formated png.
type APNG struct {
	name  string
	image apng.APNG
}

// Encode fills the background and writes the APNG.image to w in PNG format.
func (a *APNG) Encode(w io.Writer) error {
	frames := []apng.Frame{}
	for _, frame := range a.image.Frames {
		out := filledImage(frame.Image, color.White)
		img := frame
		img.Image = out
		frames = append(frames, img)
	}

	return apng.Encode(w, apng.APNG{
		Frames:    frames,
		LoopCount: 0,
	})
}

// StoreName return the file name for APNG.
func (a *APNG) StoreName() string {
	return a.name + ".png"
}
