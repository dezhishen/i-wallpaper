package main

import (
	"github.com/dezhishen/i-wallpaper/provider"
	"github.com/reujab/wallpaper"
)

type BingResponse struct {
	Images []BingImage `json:"images"`
}

type BingImage struct {
	Url string `json:"url,omitempty"`
}

func main() {
	p := provider.NewBingProvider()
	source, err := p.Random()
	check(err)
	switch source.Type {
	case provider.UrlSource:
		err = wallpaper.SetFromURL(source.Source)
		check(err)
	case provider.FileSource:
		err = wallpaper.SetFromFile(source.Source)
		check(err)
	default:
		panic("unsupport")
	}
	err = wallpaper.SetMode(wallpaper.Crop)
	check(err)
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}
