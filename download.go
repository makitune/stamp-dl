package main

import (
	"bufio"
	"context"
	"errors"
	"html"
	"image"
	"net/http"
	"strings"

	"golang.org/x/sync/errgroup"
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

	eg, ctx := errgroup.WithContext(context.TODO())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, u := range urls {
		u := u
		eg.Go(func() error {
			i, err := download(ctx, u)
			if err != nil {
				return err
			}

			s.images = append(s.images, i)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		cancel()
		return nil, err
	}
	return &s, nil
}

func download(ctx context.Context, urlString string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
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
