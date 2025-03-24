package service

import (
	"context"
	"crypto/sha256"
	"math/big"
	"short_url/domain"
	"short_url/repository"
	"time"

	"github.com/to404hanga/pkg404/logger"
)

type CachedShortUrlService struct {
	repo    repository.ShortUrlRepository
	l       logger.Logger
	suffix  string
	weights []int
}

var _ ShortUrlService = (*CachedShortUrlService)(nil)

func NewCachedShortUrlService(repo repository.ShortUrlRepository, l logger.Logger, suffix string, weights []int) *CachedShortUrlService {
	return &CachedShortUrlService{
		repo:    repo,
		l:       l,
		suffix:  suffix,
		weights: weights,
	}
}

func (s *CachedShortUrlService) Create(ctx context.Context, originUrl string) (string, error) {
	baseSuffix := ""
	for {
		shortUrl := s.generateShortUrl(originUrl, baseSuffix)
		err := s.repo.InsertShortUrl(ctx, shortUrl, originUrl)
		switch err {
		case nil, repository.ErrUniqueIndexConflict:
			return shortUrl, nil
		case repository.ErrPrimaryKeyConflict:
			baseSuffix += s.suffix
		default:
			return "", err
		}
	}
}

func (s *CachedShortUrlService) Redirect(ctx context.Context, shortUrl string) (string, error) {
	return s.repo.GetOriginUrlByShortUrl(ctx, shortUrl)
}

func (s *CachedShortUrlService) CleanExpired(ctx context.Context) error {
	now := time.Now().Unix()
	return s.repo.CleanExpired(ctx, now)
}

func (s *CachedShortUrlService) generateShortUrl(originUrl, suffix string) string {
	hashed := sha256.Sum256([]byte(originUrl + suffix))
	bytes := hashed[:]

	// 统计前导零数量
	leadingZeros := 0
	for _, b := range bytes {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	// 处理全零的特殊情况
	if leadingZeros == len(bytes) {
		return "0"
	}

	// 将字节数组转换为大整数
	num := new(big.Int).SetBytes(bytes)
	base := big.NewInt(62)
	zero := big.NewInt(0)
	var encoded []byte

	// 通过取余生成Base62字符
	for num.Cmp(zero) > 0 {
		mod := new(big.Int)
		num.DivMod(num, base, mod)
		encoded = append(encoded, domain.BASE62CHARSET[mod.Int64()])
	}

	// 反转字符顺序
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	// 补充前导零对应的字符
	for i := 0; i < leadingZeros; i++ {
		encoded = append([]byte{'0'}, encoded...)
	}

	// 计算校验位
	shortUrl := string(encoded)[:6]
	sum := s.sum(shortUrl)
	shortUrl = shortUrl[:3] + string(domain.BASE62CHARSET[sum]) + shortUrl[3:]

	return shortUrl
}

func (s *CachedShortUrlService) CheckShortUrl(ctx context.Context, shortUrl string) bool {
	if len(shortUrl) != 7 {
		return false
	}
	expected := domain.Base62NumberTable[string(shortUrl[3])]
	shortUrl = string(shortUrl[:3] + shortUrl[4:])
	sum := s.sum(shortUrl)
	return sum == expected
}

func (s *CachedShortUrlService) sum(shortUrl6 string) int {
	sum := 0
	for i := 0; i < len(shortUrl6); i++ {
		sum += domain.Base62NumberTable[string(shortUrl6[i])] * s.weights[i]
	}
	return sum % 62
}
