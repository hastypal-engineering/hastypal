package service

import (
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"github.com/google/uuid"
	"time"
)

type CreateBookingSessionService struct {
	sessionRepository  types.Repository[types.BookingSession]
	businessRepository types.Repository[types.Business]
}

func NewCreateBookingSessionService(
	sessionRepository types.Repository[types.BookingSession],
	businessRepository types.Repository[types.Business],
) *CreateBookingSessionService {
	return &CreateBookingSessionService{
		sessionRepository:  sessionRepository,
		businessRepository: businessRepository,
	}
}

func (s *CreateBookingSessionService) Execute(chatId string, businessName string) error {
	filter := types.Filter{
		Name:    "name",
		Operand: constants.Equal,
		Value:   businessName,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	business, findOneErr := s.businessRepository.FindOne(criteria)

	if findOneErr != nil {
		return findOneErr
	}

	session := types.BookingSession{
		Id:         uuid.New().String(),
		BusinessId: business.Id,
		ChatId:     chatId,
		ServiceId:  "",
		Date:       "",
		Hour:       "",
		CreatedAt:  time.Now().String(),
		Ttl:        time.Minute.Milliseconds() * 5,
	}

	if err := s.sessionRepository.Save(session); err != nil {
		return err
	}

	return nil
}
