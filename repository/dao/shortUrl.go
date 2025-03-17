package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type GormShortUrlDAO struct {
	db *gorm.DB
}

var _ ShortUrlDAO = (*GormShortUrlDAO)(nil)

var (
	ErrPrimaryKeyConflict  = errors.New("primary key conflict")
	ErrUniqueIndexConflict = errors.New("unique index conflict")
)

func NewGormShortUrlDAO(db *gorm.DB) ShortUrlDAO {
	return &GormShortUrlDAO{db: db}
}

func (g *GormShortUrlDAO) tableName(shortUrl string) string {
	return fmt.Sprintf("_%s", string(shortUrl[0]))
}

func (g *GormShortUrlDAO) Insert(ctx context.Context, su ShortUrl) error {
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txt := tx.WithContext(ctx).Table(g.tableName(su.ShortUrl))
		if err := txt.Create(&su).Error; err != nil {
			var shortUrl ShortUrl
			if err = txt.Where("short_url = ?", su.ShortUrl).Find(&shortUrl).Error; err != nil {
				return err
			}
			if shortUrl.OriginUrl == su.OriginUrl {
				return ErrPrimaryKeyConflict
			} else {
				return ErrUniqueIndexConflict
			}
		}
		return nil
	})
	return err
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
		group   errgroup.Group
		lock    sync.Mutex
	)
	for i := 0; i < 62; i++ {
		group.Go(func() error {
			tableName := "short_url_" + string(base62[i])
			for {
				var ret []string
				// 查询可删除列表
				err := g.db.WithContext(ctx).Table(tableName).Select("short_url").
					Where("expired_at < ?", now).Order("expired_at ASC").Limit(100).
					Find(&ret).Error
				if err != nil {
					return err
				}
				if len(ret) == 0 {
					break // 无更多数据可删除
				}
				err = g.db.WithContext(ctx).Table(tableName).Where("short_url IN ?", ret).Delete(&ShortUrl{}).Error
				if err != nil {
					return err
				}

				lock.Lock()
				retList = append(retList, ret...)
				lock.Unlock()

				time.Sleep(100 * time.Millisecond) // 避免高频操作压垮数据库
			}
			return nil
		})
	}
	return retList, group.Wait()
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
