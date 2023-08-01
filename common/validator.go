package common

import (
	"fmt"
	"task/internal/domain/Errors"
)

func ValidateCurrency(trCurrency string, acCurrency string, amount float64) (float64, error) {
	const op = "common.ValidateCurrency"

	switch trCurrency {
	case "USD":
		switch acCurrency {
		case "USD":
			return amount, nil
		case "EUR":
			return amount * 0.9, nil
		case "RUB":
			return amount * 70, nil
		}
	case "EUR":
		switch acCurrency {
		case "USD":
			return amount * 1.1, nil
		case "EUR":
			return amount, nil
		case "RUB":
			return amount * 80, nil
		}
	case "RUB":
		switch acCurrency {
		case "USD":
			return amount * 0.014, nil
		case "EUR":
			return amount * 0.013, nil
		case "RUB":
			return amount, nil
		}
	}
	return 0, fmt.Errorf("%s: %w", op, Errors.ErrInvalidCurrency)

}
