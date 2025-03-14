package ioc

import (
	"fmt"
	mysharding "short_url/pkg/sharding"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/sharding"
)

func InitDB() *gorm.DB {
	type Config struct {
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		Database    string `yaml:"database"`
		TablePrefix string `yaml:"tablePrefix"`
	}
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   cfg.TablePrefix,
		},
	})
	if err != nil {
		panic(err)
	}
	db.Use(sharding.Register(sharding.Config{
		ShardingKey:       "short_url",
		NumberOfShards:    62,
		ShardingAlgorithm: mysharding.CustomShardingAlgorithm,
		ShardingSuffixs: func() (suffixs []string) {
			return []string{
				"_0", "_1", "_2", "_3", "_4", "_5", "_6", "_7", "_8", "_9",
				"_A", "_B", "_C", "_D", "_E", "_F", "_G", "_H", "_I", "_J", "_K", "_L", "_M",
				"_N", "_O", "_P", "_Q", "_R", "_S", "_T", "_U", "_V", "_W", "_X", "_Y", "_Z",
				"_a", "_b", "_c", "_d", "_e", "_f", "_g", "_h", "_i", "_j", "_k", "_l", "_m",
				"_n", "_o", "_p", "_q", "_r", "_s", "_t", "_u", "_v", "_w", "_x", "_y", "_z",
			}
		},
	}, "short_url"))

	// dao.InitTables(db)

	return db
}
