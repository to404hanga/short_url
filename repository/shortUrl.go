package repository

import (
	"context"
	"short_url/repository/cache"
	"short_url/repository/dao"
	"time"

	"github.com/to404hanga/pkg404/cachex/lru"
	"github.com/to404hanga/pkg404/logger"
)

type CachedShortUrlRepository struct {
	lru           *lru.Cache
	lruExpiration time.Duration
	cache         cache.ShortUrlCache
	dao           dao.ShortUrlDAO
	l             logger.Logger
}

type lruItem struct {
	originUrl string
	expiredAt int64
}

var _ ShortUrlRepository = (*CachedShortUrlRepository)(nil)

var (
	ErrPrimaryKeyConflict  = dao.ErrPrimaryKeyConflict
	ErrUniqueIndexConflict = dao.ErrUniqueIndexConflict
)

func NewCachedShortUrlRepository(lruSize int, lruExpiration time.Duration, cache cache.ShortUrlCache, dao dao.ShortUrlDAO, l logger.Logger) *CachedShortUrlRepository {
	lru, err := lru.New(lruSize)
	if err != nil {
		panic(err)
	}
	return &CachedShortUrlRepository{
		lru:           lru,
		lruExpiration: lruExpiration,
		cache:         cache,
		dao:           dao,
		l:             l,
	}
}

func (c *CachedShortUrlRepository) GetOriginUrlByShortUrl(ctx context.Context, shortUrl string) (string, error) {
	now := time.Now().Unix()

	// 先查本地缓存，若本地缓存存在直接返回
	if val, ok := c.lru.Get(shortUrl); ok {
		if item, ok := val.(lruItem); ok && item.expiredAt >= now {
			return item.originUrl, nil
		}
	}

	// 再查redis，若redis存在直接返回
	if originUrl, err := c.cache.Get(ctx, shortUrl); err == nil {
		return originUrl, nil
	}

	// 最后查数据库
	su, err := c.dao.FindByShortUrlWithExpired(ctx, shortUrl, now)
	if err != nil {
		return "", err
	}

	go func() {
		// 异步更新redis缓存
		if err = c.cache.Set(ctx, shortUrl, su.OriginUrl); err != nil {
			c.l.Error("failed to set redis cache",
				logger.Error(err),
				logger.String("short_url", shortUrl),
				logger.String("origin_url", su.OriginUrl),
			)
		}
		// 异步更新本地lru缓存
		c.lru.Add(shortUrl, lruItem{
			originUrl: su.OriginUrl,
			expiredAt: int64(c.lruExpiration.Seconds()),
		})
	}()

	return su.OriginUrl, nil
}

func (c *CachedShortUrlRepository) InsertShortUrl(ctx context.Context, shortUrl, originUrl string) error {
	return c.dao.Insert(ctx, dao.ShortUrl{
		ShortUrl:  shortUrl,
		OriginUrl: originUrl,
		ExpiredAt: time.Now().AddDate(1, 0, 0).Unix(), // 有效期一年
	})
}

func (c *CachedShortUrlRepository) DeleteShortUrlByShortUrl(ctx context.Context, shortUrl string) error {
	err := c.dao.DeleteByShortUrl(ctx, shortUrl)
	if err == nil {
		go func() {
			// 异步删除redis缓存
			if err = c.cache.Del(ctx, shortUrl); err != nil {
				c.l.Error("failed to delete redis cache",
					logger.Error(err),
					logger.String("short_url", shortUrl),
				)
			}
			// 异步删除本地lru缓存
			c.lru.Remove(shortUrl)
		}()
	}
	return err
}

func (c *CachedShortUrlRepository) CleanExpired(ctx context.Context, now int64) error {
	deleteList, err := c.dao.DeleteExpiredList(ctx, now)
	if err == nil {
		go func() {
			for _, shortUrl := range deleteList {
				// 异步删除redis缓存
				if err = c.cache.Del(ctx, shortUrl); err != nil {
					c.l.Error("failed to delete redis cache",
						logger.Error(err),
						logger.String("short_url", shortUrl),
					)
				}
				// 异步删除本地lru缓存
				c.lru.Remove(shortUrl)
			}
		}()
	}
	return err
}
