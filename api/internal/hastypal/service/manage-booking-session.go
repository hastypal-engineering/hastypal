package service

import (
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
)

type ManageBookingSessionService struct {
	sessionRepository types.Repository[types.BookingSession]
}

func NewManageBookingSessionService(
	sessionRepository types.Repository[types.BookingSession],
) *ManageBookingSessionService {
	return &ManageBookingSessionService{
		sessionRepository: sessionRepository,
	}
}

func (s *ManageBookingSessionService) Execute(updatedSession types.BookingSession) error {
	filter := types.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   updatedSession.Id,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	currentSession, findOneErr := s.sessionRepository.FindOne(criteria)

	if findOneErr != nil {
		if err := s.sessionRepository.Save(updatedSession); err != nil {
			return err
		}

		return nil
	}

	mergedSession := helper.Merge[types.BookingSession](currentSession, updatedSession)

	if err := s.sessionRepository.Update(mergedSession); err != nil {
		return err
	}

	return nil
}
