package provider

import (
	"fmt"

	"github.com/dezhishen/i-wallpaper/pkg/config"
	log "github.com/sirupsen/logrus"
)

type sourceType int

const (
	FileSource sourceType = iota
	UrlSource
)

type PickType int

const (
	TodayPick PickType = iota
	RandomPick
)

var PickTypeLabel = []string{"今日", "随机"}

type ImageSource struct {
	Type   sourceType
	Source string
}

type Provider interface {
	Name() string
	GetTody() (*ImageSource, error)
	Random() (*ImageSource, error)
}

type ProviderFactory interface {
	Name() string
	New(map[string]interface{}) (Provider, error)
}

var _allFactories []ProviderFactory

func addFactory(f ProviderFactory) {
	_allFactories = append(_allFactories, f)
}

var _allProvider = make(map[string]Provider)

func _register(p Provider) {
	if p != nil {
		log.Infof("register provider: %s", p.Name())
		_allProvider[p.Name()] = p
	}
}

func Get(name string) (Provider, bool) {
	p, ok := _allProvider[name]
	if !ok {
		return nil, ok
	}
	return p, ok
}

func GetAllProviders() []string {
	var result []string
	for k := range _allProvider {
		result = append(result, k)
	}
	return result
}

func Init() error {
	for _, f := range _allFactories {
		name := f.Name()
		providerConfigMap := config.GetStringMap(fmt.Sprintf("provider.all.%s", name))
		if providerConfigMap == nil {
			providerConfigMap = make(map[string]interface{})
		}
		p, err := f.New(providerConfigMap)
		if err != nil {
			return err
		}
		_register(p)
	}
	return nil
}
