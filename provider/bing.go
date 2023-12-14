package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type bingImage struct {
	Url string `json:"url,omitempty"`
}

type bingResult struct {
	Images []bingImage `json:"images"`
}

// "https://cn.bing.com"

type bingProvider struct {
	endpoint string
}

func NewBingProvider() Provider {
	return &bingProvider{
		endpoint: "https://cn.bing.com",
	}
}

// impl provider provider
func (p *bingProvider) GetTody() (*ImageSource, error) {
	url := p.endpoint + "/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	var bingR bingResult
	err = json.Unmarshal(body, &bingR)
	if err != nil {
		return nil, err
	}
	return &ImageSource{
		Type:   UrlSource,
		Source: fmt.Sprintf("%s/%s", p.endpoint, bingR.Images[0].Url),
	}, nil
}

func (p *bingProvider) Random() (*ImageSource, error) {
	url := p.endpoint + "/HPImageArchive.aspx?format=js&idx=0&n=8&mkt=zh-CN"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var bingR bingResult
	err = json.Unmarshal(body, &bingR)
	if err != nil {
		return nil, err
	}
	return &ImageSource{
		Type:   UrlSource,
		Source: fmt.Sprintf("%s/%s", p.endpoint, pick(bingR.Images).Url),
	}, nil
}

func pick(reasons []bingImage) bingImage {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	return reasons[r.Intn(len(reasons))]
}
