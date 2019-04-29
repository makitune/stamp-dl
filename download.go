package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"html"
	"image"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/sync/errgroup"
)

// lineDataPreview is a object relational mapping structure
type lineDataPreview struct {
	Type              string `json:"type"`
	ID                string `json:"id"`
	StaticURL         string `json:"staticUrl"`
	FallbackStaticURL string `json:"fallbackStaticUrl"`
	AnimationURL      string `json:"animationUrl"`
	PopupURL          string `json:"popupUrl"`
	SoundURL          string `json:"soundUrl"`
}

func stampTypeURL(ldp *lineDataPreview) (*url.URL, error) {
	switch ldp.Type {
	case "static":
		return url.Parse(ldp.StaticURL)
	case "animation":
		return url.Parse(ldp.AnimationURL)
	case "popup":
		return url.Parse(ldp.PopupURL)
	case "sound":
		return url.Parse(ldp.SoundURL)
	case "animation_sound":
		return nil, errors.New("ボイス・サウンド付きスタンプには対応していません")
	default:
		return nil, errors.New("対応していません")
	}
}

// lineDataPreviews is a collection object for lineDataPreview
type lineDataPreviews struct {
	Title        string
	DataPreviews []*lineDataPreview
}

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
	var dataPreviews []*lineDataPreview

	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mdCMN38Item01Ttl") {
			start := strings.Index(line, ">")
			end := strings.Index(line, "</p>")
			title = line[start+1 : end]
		}

		if strings.Contains(line, "data-preview") {
			start := strings.Index(line, "{")
			end := strings.Index(line, "}")
			j := html.UnescapeString(line[start : end+1])

			pd := new(lineDataPreview)
			if err := json.Unmarshal([]byte(j), pd); err != nil {
				return nil, err
			}
			dataPreviews = append(dataPreviews, pd)
		}
	}

	ldp := &lineDataPreviews{
		Title:        html.UnescapeString(title),
		DataPreviews: dataPreviews[1:],
	}

	var stickers []LineSticker
	eg, ctx := errgroup.WithContext(context.TODO())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, dp := range ldp.DataPreviews {
		id := dp.ID
		u, err := stampTypeURL(dp)
		if err != nil {
			return nil, err
		}

		eg.Go(func() error {
			i, err := download(ctx, u.String())
			if err != nil {
				return err
			}

			stickers = append(stickers, LineSticker{
				ID:    id,
				Image: i,
			})
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		cancel()
		return nil, err
	}

	return &LineStamp{
		Title:    ldp.Title,
		Stickers: stickers,
	}, nil
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
