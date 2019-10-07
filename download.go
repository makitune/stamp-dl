package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
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

func (p *lineDataPreview) URL() (*url.URL, error) {
	typ, err := ParseStickerType(p.Type)
	if err != nil {
		return nil, err
	}
	switch typ {
	case LineStickerStatic:
		return url.Parse(p.StaticURL)
	case LineStickerAnimation:
		return url.Parse(p.AnimationURL)
	case LineStickerCustom:
		return url.Parse(p.StaticURL)
	default:
		return nil, fmt.Errorf("%sには対応していません", typ.Name())
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
	var stickers []Encoder
	eg, ctx := errgroup.WithContext(context.TODO())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, dp := range dps.dataPreviews {
		d := dp

		eg.Go(func() error {
			s, err := download(ctx, d)
			if err != nil {
				return err
			}

			stickers = append(stickers, s)
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

func download(ctx context.Context, ldp *lineDataPreview) (*FileEncoder, error) {
	u, err := ldp.URL()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	typ, err := ParseStickerType(ldp.Type)
	if err != nil {
		return nil, err
	}
	img, err := typ.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return &FileEncoder{name: ldp.ID + ".png", img: img}, nil
}
