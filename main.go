package main

import (
	"github.com/dezhishen/i-wallpaper/pkg/config"
	"github.com/dezhishen/i-wallpaper/pkg/font"
	"github.com/dezhishen/i-wallpaper/pkg/gui"
	"github.com/dezhishen/i-wallpaper/pkg/provider"
)

func main() {
	// 加载配置
	err := config.Init("config", "config")
	if err != nil {
		panic(err)
	}
	// 初始化字体
	defer font.BeforeDestory()
	font.Init()
	err = provider.Init()
	if err != nil {
		panic(err)
	}
	err = gui.Start()
	if err != nil {
		panic(err)
	}
}
