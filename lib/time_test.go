package lib

import (
	"testing"
	"time"
)

func TestToYYYMMDD(t *testing.T) {
	val := time.Date(2022, 3, 14, 0,0,0,0,time.UTC)
	expectedVal := "2022-03-14"

	if retVal := ToYYYMMDD(val); retVal != expectedVal {
		t.Fatal("Failed to parse", retVal)
	}
}
