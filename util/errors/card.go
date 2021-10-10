package errors

import (
	"fmt"
	"net/http"

	"gitlab.com/renodesper/spenmo-test/util/error"
)

var (
	InvalidCardNumber  = error.NewError(http.StatusBadRequest, "CA1001", fmt.Errorf("invalid card number"))
	InvalidExpiryMonth = error.NewError(http.StatusBadRequest, "CA1001", fmt.Errorf("invalid expiry month"))
	InvalidExpiryYear  = error.NewError(http.StatusBadRequest, "CA1002", fmt.Errorf("invalid expiry year"))
	ExpiredCard        = error.NewError(http.StatusBadRequest, "CA1003", fmt.Errorf("expired card"))
	InvalidCVV         = error.NewError(http.StatusBadRequest, "CA1004", fmt.Errorf("invalid cvv"))
	FailedCardsDelete  = error.NewError(http.StatusInternalServerError, "CA1005", fmt.Errorf("failed to delete cards"))
	InvalidCardOwner   = error.NewError(http.StatusBadRequest, "CA1006", fmt.Errorf("invalid card owner"))
)
