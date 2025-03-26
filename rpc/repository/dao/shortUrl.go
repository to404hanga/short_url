package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"short_url/rpc/domain"
	"sync"
	"time"

	"github.com/to404hanga/pkg404/logger"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type GormShortUrlDAO struct {
	db *gorm.DB
	l  logger.Logger
}

var _ ShortUrlDAO = (*GormShortUrlDAO)(nil)

var (
	ErrPrimaryKeyConflict  = errors.New("primary key conflict")
	ErrUniqueIndexConflict = errors.New("unique index conflict")
	ErrDataNotFound        = gorm.ErrRecordNotFound
)

func NewGormShortUrlDAO(db *gorm.DB, l logger.Logger) ShortUrlDAO {
	return &GormShortUrlDAO{
		db: db,
		l:  l,
	}
}

func (g *GormShortUrlDAO) tableName(shortUrlOrSuffix string) string {
	if len(shortUrlOrSuffix) == 1 {
		return "short_url_" + shortUrlOrSuffix
	}
	return fmt.Sprintf("short_url_%s", string(shortUrlOrSuffix[0]))
}

func (g *GormShortUrlDAO) Insert(ctx context.Context, su ShortUrl) error {
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if tx.WithContext(ctx).Table(g.tableName(su.ShortUrl)).Create(&su).Error != nil {
			var shortUrl ShortUrl
			if err := tx.WithContext(ctx).Table(g.tableName(su.ShortUrl)).Where("short_url = ?", su.ShortUrl).Find(&shortUrl).Error; err != nil {
				return err
			}
			if shortUrl.OriginUrl != su.OriginUrl {
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
	newCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(62)
	for i := 0; i < 62; i++ {
		go func(internalCtx context.Context, suffix string) {
			defer wg.Done()
			select {
			case <-internalCtx.Done():
				return
			default:
				var internalSu ShortUrl
				if err := g.db.WithContext(internalCtx).
					Table(g.tableName(suffix)).
					Where("origin_url = ?", originUrl).
					Where("expired_at > ?", now).
					First(&internalSu).Error; err == nil {
					su = internalSu
					cancel()
				}
			}
		}(newCtx, string(domain.BASE62CHARSET[i]))
	}
	wg.Wait()
	if su.ShortUrl == "" {
		return ShortUrl{}, ErrDataNotFound
	}
	return su, nil
}

func (g *GormShortUrlDAO) FindByOriginUrlWithExpiredV1(ctx context.Context, originUrl string, now int64) (ShortUrl, error) {
	var (
		su   ShortUrl
		lock sync.Mutex
	)
	g.executeUnshardedQuery(ctx, func(iCtx context.Context, suffix string, db *gorm.DB) error {
		var internalSu ShortUrl
		if err := db.WithContext(iCtx).
			Table(g.tableName(suffix)).
			Where("origin_url =?", originUrl).
			Where("expired_at >?", now).
			First(&internalSu).Error; err != nil {
			g.l.Error("FindByOriginUrlWithExpiredV1 failed",
				logger.Error(err),
				logger.String("suffix", suffix),
				logger.String("origin_url", originUrl),
				logger.Int64("expired_at", now),
			)
			return err
		}
		lock.Lock()
		su = internalSu
		lock.Unlock()
		return nil
	})
	if su.ShortUrl == "" {
		return ShortUrl{}, ErrDataNotFound
	}
	return su, nil
}

func (g *GormShortUrlDAO) FindByOriginUrl(ctx context.Context, originUrl string) (ShortUrl, error) {
	var su ShortUrl
	newCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(62)
	for i := 0; i < 62; i++ {
		go func(internalCtx context.Context, suffix string) {
			defer wg.Done()
			select {
			case <-internalCtx.Done():
				return
			default:
				var internalSu ShortUrl
				if err := g.db.WithContext(internalCtx).
					Table(g.tableName(suffix)).
					Where("origin_url = ?", originUrl).
					First(&internalSu).Error; err == nil {
					su = internalSu
					cancel()
				}
			}
		}(newCtx, string(domain.BASE62CHARSET[i]))
	}
	wg.Wait()
	if su.ShortUrl == "" {
		return ShortUrl{}, ErrDataNotFound
	}
	return su, nil
}

func (g *GormShortUrlDAO) FindByOriginUrlV1(ctx context.Context, originUrl string) (ShortUrl, error) {
	var (
		su   ShortUrl
		lock sync.Mutex
	)
	g.executeUnshardedQuery(ctx, func(iCtx context.Context, suffix string, db *gorm.DB) error {
		var internalSu ShortUrl
		if err := db.WithContext(iCtx).
			Table(g.tableName(suffix)).
			Where("origin_url =?", originUrl).
			First(&internalSu).Error; err != nil {
			g.l.Error("FindByOriginUrlV1 failed",
				logger.Error(err),
				logger.String("suffix", suffix),
				logger.String("origin_url", originUrl),
			)
			return err
		}
		lock.Lock()
		su = internalSu
		lock.Unlock()
		return nil
	})
	if su.ShortUrl == "" {
		return ShortUrl{}, ErrDataNotFound
	}
	return su, nil
}

func (g *GormShortUrlDAO) FindExpiredList(ctx context.Context, now int64) ([]ShortUrl, error) {
	var (
		sus  []ShortUrl
		wg   sync.WaitGroup
		lock sync.Mutex
	)
	newCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	wg.Add(62)
	for i := 0; i < 62; i++ {
		go func(internalCtx context.Context, suffix string) {
			defer wg.Done()
			select {
			case <-internalCtx.Done():
				return
			default:
				var internalSus []ShortUrl
				if err := g.db.WithContext(internalCtx).
					Table(g.tableName(suffix)).
					Where("expired_at <=?", now).
					Find(&internalSus).Error; err == nil {
					lock.Lock()
					sus = append(sus, internalSus...)
					lock.Unlock()
					cancel()
				}
			}
		}(newCtx, string(domain.BASE62CHARSET[i]))
	}
	wg.Wait()
	if len(sus) == 0 {
		return nil, ErrDataNotFound
	}
	return sus, nil
}

func (g *GormShortUrlDAO) FindExpiredListV1(ctx context.Context, now int64) ([]ShortUrl, error) {
	var (
		sus  []ShortUrl
		lock sync.Mutex
	)
	g.executeUnshardedQuery(ctx, func(iCtx context.Context, suffix string, db *gorm.DB) error {
		var internalSus []ShortUrl
		err := db.WithContext(iCtx).
			Table(g.tableName(suffix)).
			Where("expired_at <=?", now).
			Find(&internalSus).Error
		if err != nil {
			g.l.Error("FindExpiredListV1 failed",
				logger.Error(err),
				logger.String("suffix", suffix),
				logger.Int64("expired_at", now),
			)
			return err
		}
		lock.Lock()
		sus = append(sus, internalSus...)
		lock.Unlock()
		return nil
	})
	if len(sus) == 0 {
		return nil, ErrDataNotFound
	}
	return sus, nil
}

// 批量执行不分表操作的抽象方法
func (g *GormShortUrlDAO) executeUnshardedQuery(ctx context.Context, fn func(iCtx context.Context, suffix string, db *gorm.DB) error) {
	var wg sync.WaitGroup
	newCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	wg.Add(62)
	for i := 0; i < 62; i++ {
		go func(internalCtx context.Context, suffix string) {
			defer wg.Done()
			select {
			case <-internalCtx.Done():
				return
			default:
				if err := fn(internalCtx, suffix, g.db); err == nil {
					cancel()
				}
			}
		}(newCtx, string(domain.BASE62CHARSET[i]))
	}
	wg.Wait()
}

func (g *GormShortUrlDAO) DeleteByShortUrl(ctx context.Context, shortUrl string) error {
	return g.db.WithContext(ctx).Table(g.tableName(shortUrl)).Where("short_url = ?", shortUrl).Delete(&ShortUrl{}).Error
}

func (g *GormShortUrlDAO) DeleteExpiredList(ctx context.Context, now int64) ([]string, error) {
	var (
		retList []string
		group   errgroup.Group
		lock    sync.Mutex
	)
	for i := 0; i < 62; i++ {
		group.Go(func() error {
			tableName := "short_url_" + string(domain.BASE62CHARSET[i])
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
