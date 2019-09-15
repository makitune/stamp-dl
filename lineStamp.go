package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"strings"

	"github.com/kettek/apng"
)

var lineStore = "https://store.line.me/stickershop"

type LineStickerType int

const (
	LineStickerStatic LineStickerType = iota
	LineStickerAnimation
	LineStickerAnimationSound
	LineStickerPopup
	LineStickerSound
	LineStickerCustom
	LineStickerUnkown
)

// LineSticker is an object for a Line stamp image
type LineSticker struct {
	ID    string
	Image StickerImage
}

type StickerImage struct {
	Type LineStickerType
	raw  interface{}
}

// StoreName is the sticker name for saving
func (s *LineSticker) StoreName() string {
	return s.ID + ".png"
}

// FilledBackgroundImage is the sticker image for saving
func (s *LineSticker) FilledBackgroundImage(clr color.Color) (*StickerImage, error) {
	switch s.Image.Type {
	case LineStickerStatic, LineStickerCustom:
		si := fillImageBackground(s.Image, color.RGBA{255, 255, 255, 255})
		return &si, nil
	case LineStickerAnimation:
		si := fillAPNGBackground(s.Image, color.RGBA{255, 255, 255, 255})
		return &si, nil
	}

	return nil, errors.New("対応していません")
}

func (si *StickerImage) Encode(w io.Writer) error {
	switch si.Type {
	case LineStickerStatic, LineStickerCustom:
		return png.Encode(w, si.raw.(image.Image))
	case LineStickerAnimation:
		return apng.Encode(w, si.raw.(apng.APNG))
	}

	return errors.New("対応していません")
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
