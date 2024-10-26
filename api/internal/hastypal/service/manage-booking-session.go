package service

import (
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"time"
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
	businessFilter := types.Filter{
		Name:    "business_id",
		Operand: constants.Equal,
		Value:   updatedSession.BusinessId,
	}

	chatFilter := types.Filter{
		Name:    "chat_id",
		Operand: constants.Equal,
		Value:   updatedSession.ChatId,
	}

	sessionTimestamp, timeParseErr := time.Parse(time.DateTime, updatedSession.CreatedAt)

	if timeParseErr != nil {
		return types.ApiError{
			Msg:      timeParseErr.Error(),
			Function: "time.Parse()",
			File:     "service/manage-booking-session.go",
			Values:   []string{updatedSession.CreatedAt},
		}
	}

	createdAtLessThanOrEqualFilter := types.Filter{
		Name:    "created_at",
		Operand: constants.LessThanOrEqual,
		Value:   sessionTimestamp.Add(time.Duration(updatedSession.Ttl)).String(),
	}

	createdAtGreaterThanOrEqualFilter := types.Filter{
		Name:    "created_at",
		Operand: constants.GreaterThanOrEqual,
		Value:   updatedSession.CreatedAt,
	}

	// created -> 12:25, valid until -> 12:30, now is -> 12:26

	// created -> 12:25, valid until -> 12:30, now is -> 12:31

	criteria := types.Criteria{Filters: []types.Filter{
		businessFilter,
		chatFilter,
		createdAtLessThanOrEqualFilter,
		createdAtGreaterThanOrEqualFilter,
	}}

	currentSession, findOneErr := s.sessionRepository.FindOne(criteria)

	if findOneErr != nil {
		if err := s.sessionRepository.Save(updatedSession); err != nil {
			return err
		}

		return nil
	}

	reflection := helper.NewReflectionHelper[types.BookingSession]()

	mergedSession := reflection.Merge(currentSession, updatedSession)

	if err := s.sessionRepository.Update(mergedSession); err != nil {
		return err
	}

	return nil
}
