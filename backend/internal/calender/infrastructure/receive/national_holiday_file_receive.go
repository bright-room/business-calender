package receive

import (
	"net.bright-room.dev/calender-api/internal/calender/application/repository"
	"net.bright-room.dev/calender-api/internal/calender/domain/model"
)

type nationalHolidayFileReceive struct{}

func (r nationalHolidayFileReceive) Receive() model.Holidays {

}

func NewNationalHolidayFileReceive() repository.NationalHolidayReceiveRepository {
	return nationalHolidayFileReceive{}
}
