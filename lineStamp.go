package main

import (
	"image"
	"image/color"
	"strings"
)

var lineStore = "https://store.line.me/stickershop"

// LineSticker is an object for a Line stamp image
type LineSticker struct {
	ID    string
	Image image.Image
}

// StoreName is the sticker name for saving
func (s *LineSticker) StoreName() string {
	return s.ID + ".png"
}

// FilledBackgroundImage is the sticker image for saving
func (s *LineSticker) FilledBackgroundImage(clr color.Color) image.Image {
	return fillBackground(s.Image, color.RGBA{255, 255, 255, 255})
}

// LineStamp is a collection object for LineSticker
type LineStamp struct {
	Title    string
	Stickers []LineSticker
}

// IsLineStoreURL returns a boolean indicating whether the string is a LINE STORE stickershop url
func IsLineStoreURL(str string) bool {
	return strings.HasPrefix(str, lineStore)
}
