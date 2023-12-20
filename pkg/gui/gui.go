package gui

import (
	"errors"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/dezhishen/i-wallpaper/pkg/apply"
	"github.com/dezhishen/i-wallpaper/pkg/config"
	"github.com/dezhishen/i-wallpaper/pkg/provider"
	"github.com/sirupsen/logrus"
)

type providerConfig struct {
	Current string `yaml:"current"`
}

type TabItemFactory interface {
	NewTabItem() *container.TabItem
}

func Start() error {
	myApp := app.New()
	myWindow := myApp.NewWindow("i壁纸")
	conf := providerConfig{}
	err := config.Unmarshal("provider", &conf)
	if err != nil {
		return err
	}
	if conf.Current == "" {
		conf.Current = "bing"
	}
	providers := provider.GetAllProviders()
	logrus.Infof("当前总共提供方式有%d种", len(providers))
	tabs := container.NewAppTabs()
	for i, v := range providers {
		logrus.Debugf("壁纸提供方式[%d]:%s", i+1, v)
		p, ok := provider.Get(v)
		if !ok {
			return errors.New("未找到指定的壁纸提供方式【" + conf.Current + "】")
		}
		if _, ok := p.(TabItemFactory); ok {
			// do something...
			tabItem := p.(TabItemFactory).NewTabItem()
			tabs.Append(tabItem)
		} else {
			logrus.Warnf("壁纸提供方式【%s】未提供配置界面", p.Name())
		}
	}
	currentProvider, _ := provider.Get(conf.Current)
	onSelectedProvider := func(e *container.TabItem) {
		currentProvider, _ = provider.Get(e.Text)
	}
	tabs.OnSelected = onSelectedProvider
	tabs.SetTabLocation(container.TabLocationLeading)
	content := container.New(layout.NewVBoxLayout())
	content.Add(tabs)
	currentPick := provider.TodayPick
	pickSelected := widget.NewSelect(
		provider.PickTypeLabel,
		func(s string) {
			switch s {
			case "随机":
				currentPick = provider.RandomPick
				return
			default:
				currentPick = provider.TodayPick
				return
			}
		},
	)
	pickSelected.SetSelected(provider.PickTypeLabel[currentPick])
	content.Add(container.New(layout.NewFormLayout(),
		widget.NewLabel("壁纸选择方式"),
		pickSelected,
	))
	content.Add(container.NewHBox(
		widget.NewButton(
			"应用",
			func() {
				logrus.Infof("apply %d of %s", currentPick, currentProvider.Name())
				apply.Apply(currentPick, currentProvider)
			},
		),
	))
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
	return nil
}
