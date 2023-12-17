package config

import (
	"errors"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

var lock = &sync.RWMutex{}

var appConfig *viper.Viper

func Init(path string, name string) error {
	lock.Lock()
	defer lock.Unlock()
	appConfig = viper.New()
	appConfig.AddConfigPath(path)
	appConfig.SetConfigName(name)
	err := appConfig.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func Unmarshal(key string, v interface{}) error {
	pv := reflect.ValueOf(v)
	if pv.Kind() != reflect.Ptr {
		return errors.New("must be a ptr")
	}
	lock.RLock()
	defer lock.RUnlock()
	return appConfig.UnmarshalKey(key, v)
}

func GetStringMap(key string) map[string]interface{} {
	lock.RLock()
	defer lock.RUnlock()
	if ok := appConfig.IsSet(key); !ok {
		return nil
	}
	return appConfig.GetStringMap(key)
}
