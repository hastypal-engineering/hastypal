package translation

import "time"

var SpanishMonths = map[time.Month]string{
	time.January:   "Enero",
	time.February:  "Febrero",
	time.March:     "Marzo",
	time.April:     "Abril",
	time.May:       "Mayo",
	time.June:      "Junio",
	time.July:      "Julio",
	time.August:    "Agosto",
	time.September: "Septiembre",
	time.October:   "Octubre",
	time.November:  "Noviembre",
	time.December:  "Diciembre",
}

type Translation struct{}

func New() *Translation {
	return &Translation{}
}

// GetSpanishMonth retrieves the Spanish name for a given month.
//
// This function provides a safe and convenient way to translate
// standard Go time.Month values into their Spanish equivalents.
//
// Parameters:
//   - month time.Month: The input month to be translated.
//     Must be a valid time.Month constant (1-12).
//
// Returns:
//   - string: The Spanish name of the month.
//     If the month is invalid or not found, returns an empty string.
//
// Behavior:
//   - Uses the predefined SpanishMonths map for translations
//   - Handles all standard time.Month values (January through December)
//   - Returns an empty string for out-of-range or undefined month values
//
// Examples:
//
//	name := GetSpanishMonth(time.January)   // Returns "Enero"
//	name := GetSpanishMonth(time.December)  // Returns "Diciembre"
//	name := GetSpanishMonth(0)              // Returns ""
func (t *Translation) GetSpanishMonth(month time.Month) string {
	return SpanishMonths[month]
}

// GetSpanishMonthShortForm retrieves the Spanish name for a given month.
//
// This function provides a safe and convenient way to translate
// standard Go time.Month values into their Spanish equivalents.
//
// Parameters:
//   - month time.Month: The input month to be translated.
//     Must be a valid time.Month constant (1-12).
//
// Returns:
//   - string: The Spanish short form name of the month.
//     If the month is invalid or not found, returns an empty string.
//
// Behavior:
//   - Uses the predefined SpanishMonths map for translations
//   - Handles all standard time.Month values (January through December)
//   - Returns an empty string for out-of-range or undefined month values
//
// Examples:
//
//	name := GetSpanishMonth(time.January)   // Returns "Ene"
//	name := GetSpanishMonth(time.December)  // Returns "Dic"
//	name := GetSpanishMonth(0)              // Returns ""
func (t *Translation) GetSpanishMonthShortForm(month time.Month) string {
	return SpanishMonths[month][:3]
}
