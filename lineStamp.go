package main

import (
	"strings"
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

// LineStamp is a collection object for LineSticker
type LineStamp struct {
	Title    string
	Stickers []Encoder
}

// IsLineStoreURL returns a boolean indicating whether the string is a LINE STORE stickershop url
func IsLineStoreURL(str string) bool {
	return strings.HasPrefix(str, lineStore)
}
