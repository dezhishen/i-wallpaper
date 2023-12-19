package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type CustomConfig struct {
	Endpoint  string `yaml:"endpoint"`
	RandCount string `yaml:"rand-count"`
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
	randCount string
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
			c.RandCount = readCountStr
		} else {
			c.RandCount = "7"
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
	if conf.RandCount == "" {
		conf.RandCount = "7"
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
	content := container.New(layout.NewVBoxLayout())
	// endpoint
	endpointLabel := widget.NewLabel("端点")
	endpointInput := widget.NewEntry()
	endpointInput.SetPlaceHolder("输入端点")
	endpointInput.SetText(p.endpoint)
	content.Add(container.New(layout.NewFormLayout(), endpointLabel, endpointInput))
	// content.Add(endpointInput)
	// content.Add(layout.NewSpacer())
	// randCount
	randCountLabel := widget.NewLabel("随机图片集数量")
	randCountInput := widget.NewSelect([]string{"7", "14", "30"}, func(s string) {
	})
	randCountInput.SetSelected(p.randCount)
	content.Add(container.New(layout.NewFormLayout(), randCountLabel, randCountInput))
	// saveBtn
	btn := widget.NewButton("保存", func() {
		p.endpoint = endpointInput.Text
		p.randCount = randCountInput.Selected
	})
	content.Add(btn)
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
