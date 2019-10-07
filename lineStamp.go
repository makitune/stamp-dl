package main

import (
	"errors"
	"strings"
)

var lineStore = "https://store.line.me/stickershop"

type LineStickerType string

const (
	LineStickerStatic         LineStickerType = "static"
	LineStickerAnimation                      = "animation"
	LineStickerAnimationSound                 = "animation_sound"
	LineStickerPopup                          = "popup"
	LineStickerSound                          = "sound"
	LineStickerCustom                         = "name"
)

var stickerNames = map[LineStickerType]string{
	LineStickerStatic:         "スタンプ", // TODO
	LineStickerAnimation:      "アニメーションスタンプ",
	LineStickerAnimationSound: "ボイス・サウンド付きスタンプ",
	LineStickerPopup:          "ポップアップスタンプ",
	LineStickerSound:          "ボイス・サウンド付きスタンプ",
	LineStickerCustom:         "カスタムスタンプ",
}

func ParseStickerType(s string) (LineStickerType, error) {
	t := LineStickerType(s)
	if _, ok := stickerNames[t]; !ok {
		return "", errors.New("対応していません")
	}
	return t, nil
}

func (t LineStickerType) Name() string {
	return stickerNames[t]
}

// LineStamp is a collection object for LineSticker
type LineStamp struct {
	Title    string
	Stickers []Encoder
}

// IsLineStoreURL returns a boolean indicating whether the string is a LINE STORE stickershop url
func IsLineStoreURL(str string) bool {
	return strings.HasPrefix(str, lineStore)
}
