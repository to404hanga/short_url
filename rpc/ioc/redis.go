package ioc

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
	cfg := Config{
		Host: "localhost",
		Port: 6379,
	}
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		panic(err)
	}

	cmd := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})
	return cmd
}
