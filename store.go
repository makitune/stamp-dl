package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// StoreStamp is a function that saves an input stamp to input dir
func StoreStamp(s *LineStamp, dir string) error {
	info, err := os.Stat(dir)
	if err != nil && !os.IsExist(err) || !info.IsDir() {
		return errors.New(dir + " というディレクトリは存在しません。")
	}

	outDir := filepath.Join(dir, s.Title)
	info, err = os.Stat(outDir)
	if err != nil && !os.IsExist(err) || !info.IsDir() {
		_ = os.Mkdir(outDir, 0755)
	}

	for _, sticker := range s.Stickers {
		absName := filepath.Join(outDir, sticker.StoreName())
		err := writeFile(sticker.FilledBackgroundImage(color.RGBA{255, 255, 255, 255}), absName)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeFile(img image.Image, name string) error {
	info, err := os.Stat(name)
	if err == nil && !info.IsDir() {
		return errors.New(name + " が既に存在するため中断しました。")
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}
