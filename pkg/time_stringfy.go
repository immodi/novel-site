package pkg

import (
	"fmt"
	"time"
)

// GetCurrentTimeRFC3339 returns the current UTC time as an RFC3339 string
func GetCurrentTimeRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// HumanReadableDate converts an RFC3339 datetime string to a human-readable date
func GetDateFromRFCStr(rfc3339Str string) (string, error) {
	t, err := time.Parse(time.RFC3339, rfc3339Str)
	if err != nil {
		return "", err
	}
	return t.Format("1/2/2006"), nil
}

// TimeAgo returns a human-readable relative time string (e.g. "a few minutes ago")
func TimeAgo(rfc3339Str string) string {
	t, err := time.Parse(time.RFC3339, rfc3339Str)
	if err != nil {
		return "unknown"
	}

	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < 2*time.Minute:
		return "a minute ago"
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < 2*time.Hour:
		return "an hour ago"
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 48*time.Hour:
		return "yesterday"
	case duration < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	case duration < 12*30*24*time.Hour: // approx months
		months := int(duration.Hours() / (24 * 30))
		if months <= 1 {
			return "a month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(duration.Hours() / (24 * 365))
		if years <= 1 {
			return "a year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}
