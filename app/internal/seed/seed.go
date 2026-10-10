// Package seed
package seed

import (
	"time"

	"github.com/adriein/hastypal/internal/business"
)

var ServiceSeed = map[string][]*business.ServiceCatalog{
	"wellness": {
		{
			Name:        "Deep Tissue Massage",
			Description: "Relieving tension with firm pressure massage.",
			Price:       80.0,
			Currency:    "USD",
			Duration:    60 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Hot Stone Therapy",
			Description: "Warm basalt stones to relax muscles and mind.",
			Price:       100.0,
			Currency:    "USD",
			Duration:    75 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Aromatherapy Facial",
			Description: "Soothing facial with essential oils.",
			Price:       65.0,
			Currency:    "USD",
			Duration:    45 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Sauna Session",
			Description: "Detoxifying heat therapy session.",
			Price:       30.0,
			Currency:    "USD",
			Duration:    30 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Reflexology",
			Description: "Foot massage targeting pressure points.",
			Price:       50.0,
			Currency:    "USD",
			Duration:    40 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
	},
	"hair": {
		{
			Name:        "Haircut & Style",
			Description: "Precision cut and finishing style.",
			Price:       40.0,
			Currency:    "USD",
			Duration:    45 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Coloring",
			Description: "Full color application with premium dyes.",
			Price:       90.0,
			Currency:    "USD",
			Duration:    90 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Balayage",
			Description: "Hand-painted highlights for a natural look.",
			Price:       120.0,
			Currency:    "USD",
			Duration:    2 * time.Hour,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Keratin Treatment",
			Description: "Smoothing treatment for frizz-free hair.",
			Price:       150.0,
			Currency:    "USD",
			Duration:    150 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Beard Trim",
			Description: "Shaping and conditioning for the beard.",
			Price:       25.0,
			Currency:    "USD",
			Duration:    20 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
	},
	"hospitality": {
		{
			Name:        "Standard Room - Night",
			Description: "Comfortable one-night stay.",
			Price:       120.0,
			Currency:    "USD",
			Duration:    24 * time.Hour,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Deluxe Suite - Night",
			Description: "Spacious suite with city view.",
			Price:       250.0,
			Currency:    "USD",
			Duration:    24 * time.Hour,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Airport Transfer",
			Description: "Private pickup and drop-off.",
			Price:       45.0,
			Currency:    "USD",
			Duration:    45 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Spa Breakfast",
			Description: "Full buffet breakfast for two.",
			Price:       35.0,
			Currency:    "USD",
			Duration:    90 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Late Checkout",
			Description: "Extended stay until evening.",
			Price:       40.0,
			Currency:    "USD",
			Duration:    6 * time.Hour,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
	},
	"restaurant": {
		{
			Name:        "Chef's Tasting Menu",
			Description: "Multi-course culinary experience.",
			Price:       85.0,
			Currency:    "USD",
			Duration:    150 * time.Minute,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Weekend Brunch",
			Description: "Bottomless brunch special.",
			Price:       45.0,
			Currency:    "USD",
			Duration:    2 * time.Hour,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Wine Pairing Dinner",
			Description: "Dinner with curated wine selections.",
			Price:       110.0,
			Currency:    "USD",
			Duration:    3 * time.Hour,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Private Chef Service",
			Description: "In-home chef for special events.",
			Price:       300.0,
			Currency:    "USD",
			Duration:    4 * time.Hour,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
		{
			Name:        "Cooking Class",
			Description: "Hands-on class with head chef.",
			Price:       70.0,
			Currency:    "USD",
			Duration:    2 * time.Hour,
			DateAdd:     time.Now(),
			DateUpd:     time.Now(),
		},
	},
}

// timeOfDay builds a time.Time carrying only a clock time, matching the
// TIME columns used by ha_business_time_slot.
func timeOfDay(hour int, minute int) time.Time {
	return time.Date(0, time.January, 1, hour, minute, 0, 0, time.UTC)
}

var WeeklyScheduleSeed = []*business.OperatingDay{
	{
		DayOfWeek: time.Monday,
		TimeSlots: []business.TimeSlot{
			{OpenTime: timeOfDay(9, 0), CloseTime: timeOfDay(13, 0)},
			{OpenTime: timeOfDay(15, 0), CloseTime: timeOfDay(17, 0)},
		},
	},
	{
		DayOfWeek: time.Tuesday,
		TimeSlots: []business.TimeSlot{
			{OpenTime: timeOfDay(9, 0), CloseTime: timeOfDay(13, 0)},
			{OpenTime: timeOfDay(15, 0), CloseTime: timeOfDay(17, 0)},
		},
	},
	{
		DayOfWeek: time.Wednesday,
		TimeSlots: []business.TimeSlot{
			{OpenTime: timeOfDay(9, 0), CloseTime: timeOfDay(13, 0)},
			{OpenTime: timeOfDay(15, 0), CloseTime: timeOfDay(17, 0)},
		},
	},
	{
		DayOfWeek: time.Thursday,
		TimeSlots: []business.TimeSlot{
			{OpenTime: timeOfDay(9, 0), CloseTime: timeOfDay(13, 0)},
			{OpenTime: timeOfDay(15, 0), CloseTime: timeOfDay(17, 0)},
		},
	},
	{
		DayOfWeek: time.Friday,
		TimeSlots: []business.TimeSlot{
			{OpenTime: timeOfDay(9, 0), CloseTime: timeOfDay(13, 0)},
			{OpenTime: timeOfDay(15, 0), CloseTime: timeOfDay(17, 0)},
		},
	},
	{
		DayOfWeek: time.Saturday,
		TimeSlots: []business.TimeSlot{
			{OpenTime: timeOfDay(9, 0), CloseTime: timeOfDay(13, 0)},
		},
	},
	{
		DayOfWeek: time.Sunday,
		IsClosed:  true,
	},
}

var HolidaySeed = []*business.Holiday{
	{
		Name:        "Thanksgiving",
		StartDate:   time.Date(2026, time.November, 26, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, time.November, 27, 0, 0, 0, 0, time.UTC),
		IsRecurring: false,
		IsClosed:    true,
		Type:        "public",
	},
	{
		Name:        "Christmas Day",
		StartDate:   time.Date(2026, time.December, 25, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2026, time.December, 25, 0, 0, 0, 0, time.UTC),
		IsRecurring: true,
		IsClosed:    true,
		Type:        "public",
	},
}

var OverrideSeed = []*business.ScheduleOverride{
	{
		Date: time.Date(2026, time.November, 25, 0, 0, 0, 0, time.UTC),
		TimeSlots: []business.TimeSlot{
			{OpenTime: timeOfDay(9, 0), CloseTime: timeOfDay(13, 0)},
		},
		Reason: "Day after Thanksgiving",
	},
	{
		Date: time.Date(2026, time.December, 24, 0, 0, 0, 0, time.UTC),
		TimeSlots: []business.TimeSlot{
			{OpenTime: timeOfDay(9, 0), CloseTime: timeOfDay(13, 0)},
		},
		Reason: "Christmas Eve",
	},
}

var BusinessSeed = []*business.Business{
	{
		Name:         "Serenity Wellness Spa",
		ContactPhone: "+1 555 0101",
		Email:        "contact@serenityspa.com",
		Address:      "12 Relaxation Ave",
		Country:      "USA",
		Lang:         "en",
		DateAdd:      time.Now(),
		DateUpd:      time.Now(),
	},
	{
		Name:         "Glam Hair Studio",
		ContactPhone: "+1 555 0102",
		Email:        "hello@glamhair.com",
		Address:      "45 Style Blvd",
		Country:      "USA",
		Lang:         "en",
		DateAdd:      time.Now(),
		DateUpd:      time.Now(),
	},
	{
		Name:         "Grand Heritage Hotel",
		ContactPhone: "+1 555 0103",
		Email:        "reservations@grandheritage.com",
		Address:      "1 King Street",
		Country:      "USA",
		Lang:         "en",
		DateAdd:      time.Now(),
		DateUpd:      time.Now(),
	},
	{
		Name:         "Savor & Flame Restaurant",
		ContactPhone: "+1 555 0104",
		Email:        "bookings@savorflame.com",
		Address:      "78 Gourmet Lane",
		Country:      "USA",
		Lang:         "en",
		DateAdd:      time.Now(),
		DateUpd:      time.Now(),
	},
}

// CatalogSeed holds the service catalog of every seeded business, keyed by business name
var CatalogSeed = map[string][]*business.ServiceCatalog{
	"Serenity Wellness Spa":    ServiceSeed["wellness"],
	"Glam Hair Studio":         ServiceSeed["hair"],
	"Grand Heritage Hotel":     ServiceSeed["hospitality"],
	"Savor & Flame Restaurant": ServiceSeed["restaurant"],
}

// ScheduleSeed holds the business schedule of every seeded business, keyed by business name
var ScheduleSeed = map[string]*business.BusinessSchedule{
	"Serenity Wellness Spa":    newSchedule(),
	"Glam Hair Studio":         newSchedule(),
	"Grand Heritage Hotel":     newSchedule(),
	"Savor & Flame Restaurant": newSchedule(),
}

// newSchedule builds a fresh schedule, the operating days and the overrides are
// persisted with IDs assigned by the database so each business needs its own copy
func newSchedule() *business.BusinessSchedule {
	weeklySchedule := make([]*business.OperatingDay, 0, len(WeeklyScheduleSeed))
	for _, day := range WeeklyScheduleSeed {
		weeklySchedule = append(weeklySchedule, &business.OperatingDay{
			DayOfWeek: day.DayOfWeek,
			IsClosed:  day.IsClosed,
			TimeSlots: day.TimeSlots,
		})
	}

	overrides := make([]*business.ScheduleOverride, 0, len(OverrideSeed))
	for _, override := range OverrideSeed {
		overrides = append(overrides, &business.ScheduleOverride{
			Date:      override.Date,
			IsClosed:  override.IsClosed,
			TimeSlots: override.TimeSlots,
			Reason:    override.Reason,
		})
	}

	return &business.BusinessSchedule{
		WeeklySchedule: weeklySchedule,
		Holidays:       HolidaySeed,
		Overrides:      overrides,
	}
}
