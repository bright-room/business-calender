package repository

import "net.bright-room.dev/calender-api/internal/calender/domain/model"

type NationalHolidayReceiveRepository interface {
	Receive() model.Holidays
}
