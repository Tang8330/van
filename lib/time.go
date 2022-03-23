package lib

import (
	"strings"
	"time"
)

const YYYYMMDD = "2006-01-02"

func ToYYYMMDD(t1 time.Time) string {
	return t1.Format(YYYYMMDD)
}

// Pass in 3/14/2021 receive 12/31/2021
func EndOfYear(t1 time.Time) time.Time {
	return time.Date(t1.Year(), 12, 31, 0, 0, 0, 0, t1.Location())
}

func StartOfYear(t1 time.Time) time.Time {
	return time.Date(t1.Year(), 1, 1, 0, 0, 0, 0, t1.Location())
}

type TimeFormatYYYYMMDD time.Time

func (t *TimeFormatYYYYMMDD) UnmarshalJSON(b []byte) (err error) {
	parsedTime, err := time.Parse(YYYYMMDD, strings.Replace(string(b), `"`, "", -1))
	if err != nil {
		return err
	}

	*t = TimeFormatYYYYMMDD(parsedTime)
	return nil
}
