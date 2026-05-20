package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgNumericToString(num pgtype.Numeric, def string) string {
	value, err := num.Value()
	if err != nil {
		return def
	}
	return value.(string)
}

func PgTimestampToString(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return "unknown"
	}
	return ts.Time.Format(time.RFC3339)
}

func PgTimestampToTime(ts pgtype.Timestamptz) (time.Time, error) {
	if !ts.Valid {
		return time.Time{}, nil
	}
	return ts.Time, nil
}

func TimeToNaturalLanguage(t time.Time) string {
	now := time.Now()
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else if t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day() {
		return fmt.Sprintf("today at %s", t.Format("15:04"))
	} else if t.Year() == now.Year() {
		return t.Format("Jan 2 at 15:04")
	} else {
		return t.Format("Jan 2, 2006 at 15:04")
	}
}

func PgTimestampToNaturalLanguage(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return "Unknown"
	}
	t := ts.Time
	return TimeToNaturalLanguage(t)
}

func StringToPgText(s string) pgtype.Text {
	var text pgtype.Text
	if err := text.Scan(s); err != nil {
		return pgtype.Text{String: "", Valid: false}
	}
	return text
}
