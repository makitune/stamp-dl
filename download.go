package main

import (
	"bufio"
	"errors"
	"html"
	"image"
	"net/http"
	"strings"
)

func fetchStamps(urls []string) ([]*LineStamp, error) {
	var stamps []*LineStamp
	for _, u := range urls {
		s, err := fetchStamp(u)
		if err != nil {
			return nil, err
		}

		stamps = append(stamps, s)
	}
	return stamps, nil
}

func fetchStamp(urlString string) (*LineStamp, error) {
	if !IsLineStoreURL(urlString) {
		return nil, errors.New(urlString + " はLINEスタンプページのURLではありません。")
	}

	resp, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}

	var title string
	var urls []string
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mdCMN38Item01Ttl") {
			start := strings.Index(line, ">")
			end := strings.Index(line, "</p>")
			title = line[start+1 : end]
		}

		if strings.Contains(line, "style=\"background-image") {
			start := strings.Index(line, "(")
			end := strings.Index(line, ";")
			urls = append(urls, line[start+1:end])
		}
	}

	s := LineStamp{
		title:  html.UnescapeString(title),
		images: []image.Image{},
	}
	for _, u := range urls {
		i, err := download(u)
		if err != nil {
			return nil, err
		}

		s.images = append(s.images, i)
	}

	return &s, nil
}

func download(urlString string) (image.Image, error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}
