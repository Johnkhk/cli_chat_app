package lib

import "time"

// FormatTimestamp formats a time.Time value into a human-readable string
func FormatTimestamp(t time.Time) string {
	return t.Format("15:04") // Shows time in 24-hour format (HH:MM)
	// Alternative formats:
	// return t.Format("3:04 PM") // 12-hour format with AM/PM
	// return t.Format("2006-01-02 15:04") // Include date
}
