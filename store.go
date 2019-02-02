package main

import (
	"errors"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
)

func storeStamp(s *LineStamp, dir string) error {
	info, err := os.Stat(dir)
	if err != nil && !os.IsExist(err) || !info.IsDir() {
		return errors.New(dir + " というディレクトリは存在しません。")
	}

	outDir := filepath.Join(dir, s.title)
	info, err = os.Stat(outDir)
	if err != nil && !os.IsExist(err) || !info.IsDir() {
		_ = os.Mkdir(outDir, 0755)
	}

	for i, img := range s.filledBackgroundImage(color.RGBA{255, 255, 255, 255}) {
		name := strconv.Itoa(i) + ".png"
		absName := filepath.Join(outDir, name)
		info, err := os.Stat(absName)
		if err == nil && !info.IsDir() {
			return errors.New(absName + " が既に存在するため中断しました。")
		}

		f, err := os.Create(absName)
		if err != nil {
			return err
		}

		defer f.Close()
		err = png.Encode(f, img)
		if err != nil {
			return err
		}
	}

	return nil
}
