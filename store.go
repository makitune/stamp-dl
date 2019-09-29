package main

import (
	"errors"
	"image/color"
	"io"
	"os"
	"path/filepath"
)

// Encoder support writing to a file
type Encoder interface {
	Encode(w io.Writer) error
	StoreName() string
}

// Store is a function that saves an input stamp to input dir
func (s *LineStamp) Store(dir string) error {
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
		err := writeFile(sticker, absName)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeFile(sticker Sticker, name string) error {
	img, err := sticker.FilledBackgroundImage(color.RGBA{255, 255, 255, 255})
	if err != nil {
		return err
	}

	info, err := os.Stat(name)
	if err == nil && !info.IsDir() {
		return errors.New(name + " が既に存在するため中断しました。")
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	defer f.Close()
	err = img.Encode(f)
	if err != nil {
		return err
	}

	return nil
}
