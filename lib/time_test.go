package lib

import (
	"testing"
	"time"
)

var randomDate = time.Date(2022, 3, 14, 0, 0, 0, 0, time.UTC)

func TestToYYYMMDD(t *testing.T) {
	expectedVal := "2022-03-14"
	if retVal := ToYYYMMDD(randomDate); retVal != expectedVal {
		t.Fatal("Failed to parse", retVal)
	}
}

func TestEndOfYear(t *testing.T) {
	expectedVal := "2022-12-31"
	if retVal := ToYYYMMDD(EndOfYear(randomDate)); retVal != expectedVal {
		t.Fatal("Failed to get end of year", retVal)
	}
}

func TestStartOfYear(t *testing.T) {
	expectedVal := "2022-01-01"
	if retVal := ToYYYMMDD(StartOfYear(randomDate)); retVal != expectedVal {
		t.Fatal("Failed to get end of year", retVal)
	}
}
