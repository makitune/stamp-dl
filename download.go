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
	switch lineStickerType(ldp) {
	case LineStickerStatic:
		return url.Parse(ldp.StaticURL)
	case LineStickerAnimation:
		return url.Parse(ldp.AnimationURL)
	case LineStickerPopup:
		return nil, errors.New("ポップアップスタンプには対応していません")
	case LineStickerSound:
		return nil, errors.New("ボイス・サウンド付きスタンプには対応していません")
	case LineStickerAnimationSound:
		return nil, errors.New("ボイス・サウンド付きスタンプには対応していません")
	case LineStickerCustom:
		return url.Parse(ldp.StaticURL)
	case LineStickerUnkown:
		return nil, errors.New("対応していません")
	}

	return nil, errors.New("製作者に連絡してください")
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
		d := dp

		eg.Go(func() error {
			s, err := download(ctx, d)
			if err != nil {
				return err
			}

			stickers = append(stickers, *s)
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

func download(ctx context.Context, ldp *lineDataPreview) (*LineSticker, error) {
	u, err := stampTypeURL(ldp)
	if err != nil {
		return nil, err
	}

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
	lst := lineStickerType(ldp)
	var img interface{}
	switch lst {
	case LineStickerStatic, LineStickerCustom:
		img, _, err = image.Decode(resp.Body)
	}

	if err != nil {
		return nil, err
	}

	return &LineSticker{
		ID: ldp.ID,
		Image: StickerImage{
			Type: lst,
			raw:  img,
		},
	}, nil
}

func lineStickerType(ldp *lineDataPreview) LineStickerType {
	switch ldp.Type {
	case "static":
		return LineStickerStatic
	case "animation":
		return LineStickerAnimation
	case "animation_sound":
		return LineStickerAnimationSound
	case "popup":
		return LineStickerPopup
	case "sound":
		return LineStickerSound
	case "name":
		return LineStickerCustom
	default:
		return LineStickerUnkown
	}
}
