package dao

import (
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) {
	var base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for i := 0; i < 62; i++ {
		db.Table("short_url_" + string(base62[i])).AutoMigrate(&ShortUrl{})
	}

	// 可行
	db.Table("short_url_a").Create(&ShortUrl{
		ShortUrl:  "abcdef",
		OriginUrl: "http://example.com",
		ExpiredAt: -1,
	})

	// 不可行
	db.Create(&ShortUrl{
		ShortUrl:  "bbcdef",
		OriginUrl: "http://example.com",
		ExpiredAt: -1,
	})
}
