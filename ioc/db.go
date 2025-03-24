package ioc

import (
	"fmt"
	"log"
	"short_url/repository/dao"
	"time"

	"github.com/spf13/viper"
	"github.com/to404hanga/pkg404/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/sharding"
)

func InitDB(l logger.Logger) *gorm.DB {
	type Config struct {
		User          string `yaml:"user"`
		Password      string `yaml:"password"`
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		Database      string `yaml:"database"`
		TablePrefix   string `yaml:"tablePrefix"`
		EnableDBInit  bool   `yaml:"enableDBInit"`
		SlowThreshold int64  `yaml:"slowThreshold"`
	}
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   cfg.TablePrefix,
		},
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			SlowThreshold: time.Duration(cfg.SlowThreshold) * time.Nanosecond, // 单位 ns
			LogLevel:      glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}
	db.Use(sharding.Register(sharding.Config{
		ShardingKey:    "short_url",
		NumberOfShards: 62,
		ShardingAlgorithm: func(columnValue any) (suffix string, err error) {
			key, ok := columnValue.(string)
			if !ok {
				return "", fmt.Errorf("invalid short_url")
			}
			firstChar := string(key[0])
			suffix = fmt.Sprintf("_%s", firstChar)
			return suffix, nil
		},
		ShardingSuffixs: func() (suffixs []string) {
			return []string{
				"_0", "_1", "_2", "_3", "_4", "_5", "_6", "_7", "_8", "_9",
				"_A", "_B", "_C", "_D", "_E", "_F", "_G", "_H", "_I", "_J", "_K", "_L", "_M",
				"_N", "_O", "_P", "_Q", "_R", "_S", "_T", "_U", "_V", "_W", "_X", "_Y", "_Z",
				"_a", "_b", "_c", "_d", "_e", "_f", "_g", "_h", "_i", "_j", "_k", "_l", "_m",
				"_n", "_o", "_p", "_q", "_r", "_s", "_t", "_u", "_v", "_w", "_x", "_y", "_z",
				// "0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
				// "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
				// "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
				// "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
				// "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			}
		},
	}, "short_url"))

	if cfg.EnableDBInit {
		log.Println("Starting database initialization...")
		dao.InitTables(db)
		log.Println("Database initialization completed.")
	}

	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(s string, i ...interface{}) {
	// g(s, logger.Field{Key: "args", Val: i})
	g(fmt.Sprintf(s, i...))
}
