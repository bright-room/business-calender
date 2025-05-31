package model

import "time"

type Day interface {
	Value() time.Time
}
