package utils

import "github.com/jackc/pgx/v5/pgtype"

func PgNumericToString(num pgtype.Numeric, def string) string {
	value, err := num.Value()
	if err != nil {
		return def
	}
	return value.(string)
}
