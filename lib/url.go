package lib

import (
	"fmt"
	"net/url"
)

func EncodeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	retURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)

	vals := parsedURL.Query().Encode()

	fmt.Println("vals", vals)
	return fmt.Sprintf("%s?%s", retURL, vals), nil
}