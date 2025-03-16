package service

import (
	"context"
	"crypto/sha256"
	"math/big"
	"short_url/repository"

	"github.com/to404hanga/pkg404/logger"
)

type shortUrlService struct {
	repo   repository.ShortUrlRepository
	l      logger.Logger
	suffix string
}

var _ ShortUrlService = (*shortUrlService)(nil)

func NewShortUrlService(repo repository.ShortUrlRepository, l logger.Logger, suffix string) *shortUrlService {
	return &shortUrlService{
		repo:   repo,
		l:      l,
		suffix: suffix,
	}
}

func (s *shortUrlService) Create(ctx context.Context, originUrl string) (string, error) {

}

func (s *shortUrlService) Redirect(ctx context.Context, originUrl string) (string, error) {

}

func (s *shortUrlService) CleanExpired(ctx context.Context) error {

}

func (s *shortUrlService) generateShortUrl(originUrl, suffix string) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
		encoded = append(encoded, charset[mod.Int64()])
	}

	// 反转字符顺序
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	// 补充前导零对应的字符
	for i := 0; i < leadingZeros; i++ {
		encoded = append([]byte{'0'}, encoded...)
	}

	return string(encoded)
}
