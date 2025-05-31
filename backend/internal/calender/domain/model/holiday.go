package model

import "time"

type Holiday interface {
	Value() time.Time
	Summary() string
}
