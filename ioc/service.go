package ioc

import (
	"short_url/repository"
	"short_url/service"

	"github.com/spf13/viper"
	"github.com/to404hanga/pkg404/logger"
)

func InitService(repo repository.ShortUrlRepository, l logger.Logger) service.ShortUrlService {
	type Config struct {
		Suffix  string `yaml:"suffix"`
		Weights []int  `yaml:"weights"`
	}
	cfg := &Config{
		Suffix:  "_TO404HANGA",
		Weights: []int{5, 67, 23, 71, 73, 79},
	}
	if err := viper.UnmarshalKey("short_url", &cfg); err != nil {
		panic(err)
	}
	return service.NewCachedShortUrlService(repo, l, cfg.Suffix, cfg.Weights)
}
