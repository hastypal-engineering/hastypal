package google

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/rotisserie/eris"
	"google.golang.org/api/calendar/v3"
)

type GoogleService interface {
	Authenticate(ctx context.Context, bussinesID string) (string, error)
	Authorize() error
	CalendarEvent(ctx context.Context, businessID int, date time.Time) error
}

type Service struct {
	logger *slog.Logger
	api    CalendarApi
	repo   GoogleRepository
}

func NewService(
	logger *slog.Logger,
	api CalendarApi,
	repo GoogleRepository,
) *Service {
	return &Service{
		logger: logger,
		api:    api,
		repo:   repo,
	}
}

func (s *Service) Authenticate(ctx context.Context, bussinesID string) (string, error) {
	url := s.api.GetAuthenticationURL(bussinesID)

	return url, nil
}

func (s *Service) Authorize(ctx context.Context, req string) error {
	parsedUrl, err := url.Parse(req)

	if err != nil {
		return eris.Wrap(err, "Error parsing the URL")
	}

	code := parsedUrl.Query().Get("code")
	state := parsedUrl.Query().Get("state")

	token, err := s.api.ExchangeToken(ctx, state, code)

	if err != nil {
		return eris.Wrap(err, "Error exchanging token")
	}

	if err := s.repo.Save(ctx, token); err != nil {
		return eris.Wrap(err, "Error storing the exchanged google token")
	}

	return nil
}

func (s *Service) CalendarEvent(ctx context.Context, businessID int, date time.Time) error {
	token, err := s.repo.GetByBusinessID(ctx, businessID)

	if err != nil {
		return eris.Wrap(err, "Error retrieving the token")
	}

	client, err := s.api.Client(ctx, token)

	if err != nil {
		return eris.Wrap(err, "Error getting the google calendar client")
	}

	//TODO: the end date and start date is the same I need to figure out how to do that
	event := &calendar.Event{
		Summary:     "Google I/O 2015",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A chance to hear more about Google's developer products.",
		Start: &calendar.EventDateTime{
			DateTime: date.Format(time.RFC3339),
			TimeZone: "Europe/Madrid",
		},
		End: &calendar.EventDateTime{
			DateTime: date.Format(time.RFC3339),
			TimeZone: "Europe/Madrid",
		},
		Organizer: &calendar.EventOrganizer{
			Email: "adria.claret@gmail.com",
			Self:  true,
		},
		Status: "confirmed",
	}

	_, err = client.Events.Insert("primary", event).Do()

	if err != nil {
		return eris.Wrap(err, "Error creating the event to the calendar")
	}

	return nil
}
