package dao

import (
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) {
	db.AutoMigrate(&ShortUrl{})

	// // 可行
	// log.Println("可行")
	// db.Table("short_url_a").Create(&ShortUrl{
	// 	ShortUrl:  "abcdef",
	// 	OriginUrl: "http://example.com",
	// 	ExpiredAt: -1,
	// })

	// // 不可行
	// log.Println("不可行")
	// db.Create(&ShortUrl{
	// 	ShortUrl:  "bbcdef",
	// 	OriginUrl: "http://example.com",
	// 	ExpiredAt: -1,
	// })
}
