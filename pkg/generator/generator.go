package generator

import (
	"crypto/sha256"
	"math/big"
	"short_url/rpc/domain"
)

func GenerateShortUrl(originUrl, suffix string, weights []int) string {
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
	sum := sum(shortUrl, weights)
	shortUrl = shortUrl[:3] + string(domain.BASE62CHARSET[sum]) + shortUrl[3:]

	return shortUrl
}

func CheckShortUrl(shortUrl string, weights []int) bool {
	if len(shortUrl) != 7 {
		return false
	}
	expected := domain.Base62NumberTable[string(shortUrl[3])]
	shortUrl = string(shortUrl[:3] + shortUrl[4:])
	sum := sum(shortUrl, weights)
	return sum == expected
}

func sum(shortUrl6 string, weights []int) int {
	sum := 0
	for i := 0; i < len(shortUrl6); i++ {
		sum += domain.Base62NumberTable[string(shortUrl6[i])] * weights[i]
	}
	return sum % 62
}
