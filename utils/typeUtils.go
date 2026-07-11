package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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

func PgDateToString(date pgtype.Date) string {
	if !date.Valid {
		return "Unknown"
	}
	return date.Time.Format("02-01-2006")
}

func PgDateToFormString(date pgtype.Date) string {
	if !date.Valid {
		return ""
	}
	return date.Time.Format("2006-01-02")
}

func TimeToNaturalLanguage(t time.Time) string {
	now := time.Now()
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
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

func StringToPgDate(s string) pgtype.Date {
	var date pgtype.Date
	if err := date.Scan(s); err != nil {
		return pgtype.Date{Time: time.Time{}, Valid: false}
	}
	return date
}

func StringToPgTimestamp(s string) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	if err := ts.Scan(s); err != nil {
		return pgtype.Timestamptz{Time: time.Time{}, Valid: false}
	}
	return ts
}

func StringToInt(s string, def int) int {
	if s == "" {
		return def
	}
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return def
	}
	return i
}

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	if pgErr.Code == "23503" || pgErr.Code == "23001" {
		return true
	}
	return false
}
