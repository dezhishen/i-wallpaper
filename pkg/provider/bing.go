package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type CustomConfig struct {
	Endpoint  string `yaml:"endpoint"`
	RandCount uint   `yaml:"rand-count"`
}

type bingImage struct {
	Url string `json:"url,omitempty"`
}

type bingResult struct {
	Images []bingImage `json:"images"`
}

// "https://cn.bing.com"

type bingProvider struct {
	endpoint  string
	randCount uint
}

type bingFactory struct {
}

func (f *bingFactory) Name() string {
	return "bing"
}
func (f *bingFactory) New(configOfMap map[string]interface{}) (Provider, error) {
	c := CustomConfig{}
	if v, ok := configOfMap["endpoint"]; ok {
		c.Endpoint = fmt.Sprintf("%v", v)
	}

	if v, ok := configOfMap["rand-count"]; ok {
		readCountStr := fmt.Sprintf("%v", v)
		if readCountStr != "" {
			count, err := strconv.Atoi(readCountStr)
			if err != nil {
				return nil, err
			}
			c.RandCount = uint(count)
		}
	}
	return NewBingProvider(c), nil
}
func init() {
	addFactory(&bingFactory{})
}

func NewBingProvider(conf CustomConfig) Provider {
	if conf.Endpoint == "" {
		conf.Endpoint = "https://cn.bing.com"
	}
	if conf.RandCount == 0 {
		conf.RandCount = 8
	}
	return &bingProvider{
		endpoint:  conf.Endpoint,
		randCount: conf.RandCount,
	}
}

func (p *bingProvider) Name() string {
	return "bing"
}

func (p *bingProvider) NewTabItem() *container.TabItem {
	endpointLabel := widget.NewLabel("端点")
	content := container.New(layout.NewFormLayout())
	content.Add(endpointLabel)
	input := widget.NewEntry()
	input.SetPlaceHolder("输入端点")
	input.SetText(p.endpoint)
	content.Add(input)
	return container.NewTabItem(p.Name(), content)
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
	var bingR bingResult
	err = json.Unmarshal(body, &bingR)
	if err != nil {
		fmt.Println(string(body))
		return nil, err
	}
	return &ImageSource{
		Type:   UrlSource,
		Source: fmt.Sprintf("%s/%s", p.endpoint, strings.ReplaceAll(bingR.Images[0].Url, "1920x1080", "UHD")),
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
		fmt.Println(string(body))
		return nil, err
	}
	return &ImageSource{
		Type:   UrlSource,
		Source: fmt.Sprintf("%s/%s", p.endpoint, strings.ReplaceAll(pick(bingR.Images).Url, "1920x1080", "UHD")),
	}, nil
}

func pick(reasons []bingImage) bingImage {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	return reasons[r.Intn(len(reasons))]
}
