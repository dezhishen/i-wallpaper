package apply

import (
	"errors"

	"github.com/dezhishen/i-wallpaper/pkg/provider"
	"github.com/reujab/wallpaper"
)

func Apply(pickType provider.PickType, p provider.Provider) error {
	var source *provider.ImageSource
	var err error
	switch pickType {
	case provider.TodayPick:
		source, err = p.GetTody()
	case provider.RandomPick:
		source, err = p.GetTody()
	default:
		err = errors.New("不支持的壁纸选择方式")
	}
	if err != nil {
		return err
	}
	switch source.Type {
	case provider.UrlSource:
		err = wallpaper.SetFromURL(source.Source)
	case provider.FileSource:
		err = wallpaper.SetFromFile(source.Source)
	default:
		panic("unsupport")
	}
	if err != nil {
		return err
	}
	err = wallpaper.SetMode(wallpaper.Crop)
	return err
}
