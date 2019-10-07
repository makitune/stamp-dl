package main

import (
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/kettek/apng"
)

type PNGDecoder struct{}

// Decode decodes an png imag.
func (p *PNGDecoder) DecodeFrom(r io.Reader) (*PNG, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return &PNG{image: img}, nil
}

// PNG is an object for image formated png.
type PNG struct {
	image image.Image
}

// Encode fills the background and writes the PNG.image to w in PNG format.
func (p *PNG) EncodeTo(w io.Writer) error {
	img := filledImage(p.image, color.White)
	return png.Encode(w, img)
}

type APNGDecoder struct{}

func (a *APNGDecoder) DecodeFrom(r io.Reader) (*APNG, error) {
	img, err := apng.DecodeAll(r)
	if err != nil {
		return nil, err
	}
	return &APNG{image: img}, nil
}

// APNG is an object for image formated png.
type APNG struct {
	image apng.APNG
}

// Encode fills the background and writes the APNG.image to w in PNG format.
func (a *APNG) EncodeTo(w io.Writer) error {
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
