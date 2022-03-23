package lib

import (
	"testing"
	"time"
)

func TestBooking_ApproximateRevenue(t *testing.T) {
	start := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 06, 25, 0, 0, 0, 0, time.UTC)

	b := &Booking{
		From: TimeFormatYYYYMMDD(start),
		To:   TimeFormatYYYYMMDD(end),
	}

	// 85 days
	estRev := 2*3600 + 3*1200 + 4*110
	if approxRev := b.ApproximateRevenue(110, 1200, 3600); approxRev != estRev {
		t.Fatal("Numbers do not line up", estRev, approxRev)
	}
}
