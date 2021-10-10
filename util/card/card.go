package card

import (
	"strconv"
	"time"

	"gitlab.com/renodesper/spenmo-test/util/errors"
)

func IsCardNumberValid(cardNo string) error {
	var sum int
	var alternate bool

	numberLen := len(cardNo)
	if numberLen < 13 || numberLen > 19 {
		return errors.InvalidCardNumber
	}

	for i := numberLen - 1; i > -1; i-- {
		mod, _ := strconv.Atoi(string(cardNo[i]))
		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}

		alternate = !alternate

		sum += mod
	}

	if sum%10 != 0 {
		return errors.InvalidCardNumber
	}

	return nil
}

// NOTE: Validate ExpiryMonth & ExpiryYear
func IsExpiryValid(expiryMonth, expiryYear string) error {
	currentTime := time.Now().UTC()
	currentMonth := currentTime.Month()
	currentYear := currentTime.Year()

	month, err := IsExpiryMonthValid(expiryMonth)
	if err != nil {
		return err
	}

	year, err := IsExpiryYearValid(expiryYear)
	if err != nil {
		return err
	}

	if year == currentYear && month < int(currentMonth) {
		return errors.ExpiredCard
	}

	return nil
}

func IsExpiryMonthValid(expiryMonth string) (int, error) {
	month, err := strconv.Atoi(expiryMonth)
	if err != nil || month < 1 || month > 12 {
		return month, errors.InvalidExpiryMonth
	}

	return month, nil
}

func IsExpiryYearValid(expiryYear string) (int, error) {
	currentYear := time.Now().UTC().Year()

	year, err := strconv.Atoi(expiryYear)
	if err != nil || len(expiryYear) != 4 || year < currentYear {
		return year, errors.InvalidExpiryYear
	}

	return year, nil
}

func IsCVVValid(cvv string) error {
	_, err := strconv.Atoi(cvv)
	if err != nil || len(cvv) < 3 || len(cvv) > 4 {
		return errors.InvalidCVV
	}

	return nil
}
