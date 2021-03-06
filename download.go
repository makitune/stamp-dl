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
	title        string
	dataPreviews []*lineDataPreview
}

// FetchStamps is a function that fetches Stamps from input url
func FetchStamps(urls []string) ([]*LineStamp, error) {
	var stampData []*lineDataPreviews
	for _, u := range urls {
		sd, err := fetchStampData(u)
		if err != nil {
			return nil, err
		}

		stampData = append(stampData, sd)
	}

	var stamps []*LineStamp
	for _, dps := range stampData {
		stamp, err := downloadStamp(dps)

		if err != nil {
			return nil, err
		}

		stamps = append(stamps, stamp)
	}

	return stamps, nil
}

func fetchStampData(urlString string) (*lineDataPreviews, error) {
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

	return &lineDataPreviews{
		title:        html.UnescapeString(title),
		dataPreviews: dataPreviews[1:],
	}, nil
}

func downloadStamp(dps *lineDataPreviews) (*LineStamp, error) {
	var stickers []LineSticker
	eg, ctx := errgroup.WithContext(context.TODO())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, dp := range dps.dataPreviews {
		id := dp.ID
		u, err := stampTypeURL(dp)
		if err != nil {
			return nil, err
		}

		eg.Go(func() error {
			i, err := download(ctx, u)
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
		Title:    dps.title,
		Stickers: stickers,
	}, nil
}

func download(ctx context.Context, u *url.URL) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
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
