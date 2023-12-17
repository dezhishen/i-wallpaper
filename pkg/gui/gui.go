package gui

import (
	"errors"

	"github.com/dezhishen/i-wallpaper/pkg/config"
	"github.com/dezhishen/i-wallpaper/pkg/provider"
	"github.com/sirupsen/logrus"
)

// import (
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// )

type providerConfig struct {
	Current string `yaml:"current"`
}

func Start() error {
	conf := providerConfig{}
	err := config.Unmarshal("provider", &conf)
	if err != nil {
		return err
	}
	if conf.Current == "" {
		conf.Current = "bing"
	}
	p, ok := provider.Get(conf.Current)
	if !ok {
		return errors.New("未找到指定的壁纸提供方式【" + conf.Current + "】")
	}
	providers := provider.GetAllProviders()
	logrus.Infof("当前总共提供方式有%d种", len(providers))
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		for i, v := range providers {
			logrus.Debugf("壁纸提供方式[%d]:%s", i+1, v)
		}
	}
	logrus.Infof("当前壁纸提供方式: %s", p.Name())
	return nil
}
