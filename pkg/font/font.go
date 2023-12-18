package font

import (
	"os"
	"strings"

	findfont "github.com/flopp/go-findfont"
)

func Init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//楷体:simkai.ttf
		//黑体:simhei.ttf
		if strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func BeforeDestory() {
	os.Unsetenv("FYNE_FONT")
}
