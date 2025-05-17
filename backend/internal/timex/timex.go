package timex

import (
	"time"

	"golang.org/x/xerrors"
)

const DAY = 24 * time.Hour

var JST *time.Location

func init() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	time.Local = jst
	JST = jst
}

func StartOfUnixEpoch() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, JST)
}

func NowDate() time.Time {
	now := time.Now()
	return now.Truncate(DAY)
}

type TimeRange struct {
	Begin, End time.Time
}

func (r TimeRange) DatesUntil() ([]time.Time, error) {
	if r.Begin.After(r.End) {
		return []time.Time{}, xerrors.Errorf("begin time is after end time.")
	}

	var dates []time.Time
	current := r.Begin

	for !current.After(r.End) {
		dates = append(dates, current)
		current = current.Add(DAY)
	}

	return dates, nil
}
