package dao

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type GormShortUrlDAO struct {
	db *gorm.DB
}

var _ ShortUrlDAO = (*GormShortUrlDAO)(nil)

func NewGormShortUrlDAO(db *gorm.DB) ShortUrlDAO {
	return &GormShortUrlDAO{db: db}
}

func (g *GormShortUrlDAO) tableName(shortUrl string) string {
	return fmt.Sprintf("_%s", string(shortUrl[0]))
}

func (g *GormShortUrlDAO) Insert(ctx context.Context, su ShortUrl) error {
	return g.db.WithContext(ctx).Table(g.tableName(su.ShortUrl)).Create(&su).Error
}

func (g *GormShortUrlDAO) FindByShortUrlWithExpired(ctx context.Context, shortUrl string, now int64) (ShortUrl, error) {
	var su ShortUrl
	err := g.db.WithContext(ctx).Table(g.tableName(shortUrl)).Where("short_url = ?", shortUrl).Where("expired_at > ?", now).First(&su).Error
	return su, err
}

func (g *GormShortUrlDAO) FindByShortUrl(ctx context.Context, shortUrl string) (ShortUrl, error) {
	var su ShortUrl
	err := g.db.WithContext(ctx).Table(g.tableName(shortUrl)).Where("short_url = ?", shortUrl).First(&su).Error
	return su, err
}

func (g *GormShortUrlDAO) FindByOriginUrlWithExpired(ctx context.Context, originUrl string, now int64) (ShortUrl, error) {
	var su ShortUrl
	err := g.db.WithContext(ctx).Where("origin_url =?", originUrl).Where("expired_at > ?", now).First(&su).Error
	return su, err
}

func (g *GormShortUrlDAO) FindByOriginUrl(ctx context.Context, originUrl string) (ShortUrl, error) {
	var su ShortUrl
	err := g.db.WithContext(ctx).Where("origin_url =?", originUrl).First(&su).Error
	return su, err
}

func (g *GormShortUrlDAO) FindExpiredList(ctx context.Context, now int64) ([]ShortUrl, error) {
	var sus []ShortUrl
	err := g.db.WithContext(ctx).Where("expired_at <= ?", now).Find(&sus).Error
	return sus, err
}

func (g *GormShortUrlDAO) DeleteByShortUrl(ctx context.Context, shortUrl string) error {
	return g.db.WithContext(ctx).Where("short_url = ?", shortUrl).Delete(&ShortUrl{}).Error
}

func (g *GormShortUrlDAO) DeleteExpiredList(ctx context.Context, now int64) ([]string, error) {
	base62 := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var (
		retList []string
		err     error
	)
	for i := 0; i < 62; i++ {
		if err = g.db.Transaction(func(tx *gorm.DB) error {
			var ret []string
			tableName := "short_url_" + string(base62[i])
			err := tx.WithContext(ctx).Table(tableName).Select("short_url").Where("expired_at <= ?", now).Find(&retList).Error
			if err != nil {
				return err
			}
			err = tx.WithContext(ctx).Table(tableName).Where("short_url IN ?", ret).Delete(&ShortUrl{}).Error
			if err != nil {
				return err
			}
			retList = append(retList, ret...)
			return nil
		}); err != nil {
			break
		}
	}
	return retList, err
}

func (g *GormShortUrlDAO) Transaction(ctx context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fc(tx)
	}, opts...)
}

func (g *GormShortUrlDAO) WithTransaction(ctx context.Context, fc func(txDAO ShortUrlDAO) error, opts ...*sql.TxOptions) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := &GormShortUrlDAO{db: tx}
		return fc(txDAO)
	}, opts...)
}
