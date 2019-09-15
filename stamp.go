package main

import (
	"image/color"
	"io"
)

type Stamp interface {
	Store(dir string) error
}

type Sticker interface {
	FilledBackgroundImage(clr color.Color) (StickerImage, error)
	StoreName() string
}

type StickerImage interface {
	Encode(w io.Writer) error
}
