package service

import (
	"net.bright-room.dev/calender-api/internal/calender/application/repository"
	"net.bright-room.dev/calender-api/internal/calender/domain/model"
)

type NationalHolidayReceiveService struct {
	nationalHolidayReceiveRepository repository.NationalHolidayReceiveRepository
}

func (r *NationalHolidayReceiveService) Receive() model.Holidays {
	return r.nationalHolidayReceiveRepository.Receive()
}

func NewNationalHolidayReceiveService(repository repository.NationalHolidayReceiveRepository) *NationalHolidayReceiveService {
	return &NationalHolidayReceiveService{
		nationalHolidayReceiveRepository: repository,
	}
}
