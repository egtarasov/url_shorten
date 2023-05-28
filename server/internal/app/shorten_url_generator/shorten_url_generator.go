package shorten_url_generator

import (
	"github.com/speps/go-hashids/v2"
)

// CreateShortenUrl is a function that generates a unique token from url
func CreateShortenUrl(url string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = url
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	shortUrl, _ := h.Encode([]int{45, 434, 1313})
	return shortUrl, nil
}
