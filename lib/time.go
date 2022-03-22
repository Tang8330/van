package lib

import "time"

func ToYYYMMDD(t1 time.Time) string {
	return t1.Format("2006-01-02")
}