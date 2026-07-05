package pages

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func averageUnitPrice(totalPrice pgtype.Numeric, totalQuantity int) pgtype.Numeric {
	var zero pgtype.Numeric
	_ = zero.Scan("0.00")

	if totalQuantity <= 0 {
		return zero
	}

	value, err := totalPrice.Value()
	if err != nil {
		return zero
	}

	price, ok := value.(string)
	if !ok {
		return zero
	}

	totalRat := new(big.Rat)
	if _, ok := totalRat.SetString(price); !ok {
		return zero
	}

	avgRat := new(big.Rat).Quo(totalRat, big.NewRat(int64(totalQuantity), 1))

	var result pgtype.Numeric
	if err := result.Scan(avgRat.FloatString(2)); err != nil {
		return zero
	}

	return result
}
