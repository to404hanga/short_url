package service

import (
	"context"
	"crypto/sha256"
	"math/big"
	"short_url/repository"
	"time"

	"github.com/to404hanga/pkg404/logger"
)

type shortUrlService struct {
	repo    repository.ShortUrlRepository
	l       logger.Logger
	suffix  string
	weights []int
}

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var table = map[string]int{
	"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	"a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22,
	"n": 23, "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34, "z": 35,
	"A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48,
	"N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56, "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61,
}

var _ ShortUrlService = (*shortUrlService)(nil)

func NewShortUrlService(repo repository.ShortUrlRepository, l logger.Logger, suffix string, weights []int) *shortUrlService {
	return &shortUrlService{
		repo:    repo,
		l:       l,
		suffix:  suffix,
		weights: weights,
	}
}

func (s *shortUrlService) Create(ctx context.Context, originUrl string) (string, error) {
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

func (s *shortUrlService) Redirect(ctx context.Context, shortUrl string) (string, error) {
	return s.repo.GetOriginUrlByShortUrl(ctx, shortUrl)
}

func (s *shortUrlService) CleanExpired(ctx context.Context) error {
	now := time.Now().Unix()
	return s.repo.CleanExpired(ctx, now)
}

func (s *shortUrlService) generateShortUrl(originUrl, suffix string) string {
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

	// 计算校验位
	shortUrl := string(encoded)[:6]
	sum := 0
	for i := 0; i < len(shortUrl); i++ {
		sum += table[string(shortUrl[i])] * s.weights[i]
	}
	sum %= 62
	shortUrl = shortUrl[:3] + string(charset[sum]) + shortUrl[3:]

	return shortUrl
}
