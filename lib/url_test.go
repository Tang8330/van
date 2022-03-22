package lib

import "testing"

func TestEncodeURL(t *testing.T) {
	decodedURL := "https://search.outdoorsy.com/rentals?address=San Francisco, California&raw_json=true&average_daily_pricing=true&education=true&seo_links=true"
	encodedURL := "https://search.outdoorsy.com/rentals?address=San+Francisco%2C+California&average_daily_pricing=true&education=true&raw_json=true&seo_links=true"

	retEncodedURL, err := EncodeURL(decodedURL)
	if err != nil {
		t.Fatal("Failed to encode, returned an err", err)
	}

	if retEncodedURL != encodedURL {
		t.Fatal("Failed to encode properly", retEncodedURL)
	}
}
