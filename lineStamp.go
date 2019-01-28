package main

import (
	"image"
	"image/color"
	"strings"
)

var lineStore = "https://store.line.me/stickershop"

// LineStamp is a object
type LineStamp struct {
	title  string
	images []image.Image
}

func (s *LineStamp) filledBackgroundImage(clr color.Color) []image.Image {
	var imgs []image.Image
	for _, img := range s.images {
		imgs = append(imgs, fillBackground(img, color.RGBA{255, 255, 255, 255}))
	}
	return imgs
}

// IsLineStoreURL returns a boolean indicating whether the string is a LINE STORE stickershop url
func IsLineStoreURL(str string) bool {
	return strings.HasPrefix(str, lineStore)
}
