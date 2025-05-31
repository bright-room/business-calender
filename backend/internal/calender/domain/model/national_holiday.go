package model

import (
	"fmt"
	"time"
)

type NationalHoliday struct {
	date    time.Time
	summary string
}

func (r *NationalHoliday) Value() time.Time {
	return r.date
}

func (r *NationalHoliday) Summary() string {
	return r.summary
}

func (r *NationalHoliday) ToString() string {
	return fmt.Sprintf("%+v", r)
}

func NewNationalHoliday(date time.Time, summary string) *NationalHoliday {
	return &NationalHoliday{
		date:    date,
		summary: summary,
	}
}
