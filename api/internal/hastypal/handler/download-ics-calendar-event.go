package handler

import (
	"fmt"
	"net/http"
	"time"
)

type DownloadIcsCalendarEvent struct {
}

func NewDownloadIcsCalendarEvent() *DownloadIcsCalendarEvent {
	return &DownloadIcsCalendarEvent{}
}

func (h *DownloadIcsCalendarEvent) Handler(w http.ResponseWriter, r *http.Request) error {
	icsContent := generateICS()

	// Set headers
	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", "attachment; filename=event.ics")

	// Write content
	_, err := w.Write([]byte(icsContent))
	if err != nil {
		http.Error(w, "Unable to generate .ics file", http.StatusInternalServerError)
	}

	return nil
}

func generateICS() string {
	now := time.Now()
	startTime := now.Add(24 * time.Hour).Format("20060102T150405Z")
	endTime := now.Add(25 * time.Hour).Format("20060102T150405Z")
	uid := fmt.Sprintf("%d@example.com", now.UnixNano())

	return fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Your Company//Your Product//EN
CALSCALE:GREGORIAN
BEGIN:VEVENT
UID:%s
DTSTAMP:%s
DTSTART:%s
DTEND:%s
SUMMARY:Sample Event
DESCRIPTION:This is a sample event.
LOCATION:Online
END:VEVENT
END:VCALENDAR`, uid, now.UTC().Format("20060102T150405Z"), startTime, endTime)
}
