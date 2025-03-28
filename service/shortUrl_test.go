package service

import "testing"

func TestShortUrlService_generateShortUrl(t *testing.T) {
	testCases := []struct {
		name      string
		originUrl string
		suffix    string
	}{
		{
			name:      "正常生成1",
			originUrl: "https://example.com",
			suffix:    "",
		},
		{
			name:      "正常生成2",
			originUrl: "http://123456789.com/1563sdf/set",
			suffix:    "",
		},
		{
			name:      "模拟生成冲突1",
			originUrl: "https://example.com",
			suffix:    "SHORTURL",
		},
		{
			name:      "模拟生成冲突2",
			originUrl: "http://123456789.com/1563sdf/set",
			suffix:    "SHORTURL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewCachedShortUrlService(nil, nil, "", []int{5, 67, 23, 71, 73, 79})
			got := s.generateShortUrl(tc.originUrl, tc.suffix)
			t.Log(got)
		})
	}
}
