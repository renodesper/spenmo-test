package errors

import (
	"fmt"
	"net/http"

	"gitlab.com/renodesper/spenmo-test/util/error"
)

var (
	FailedToDeleteCards           = error.NewError(http.StatusInternalServerError, "WL1000", fmt.Errorf("failed to delete cards"))
	MissingWalletOwner            = error.NewError(http.StatusBadRequest, "WL1001", fmt.Errorf("missing wallet owner"))
	InvalidWalletOwner            = error.NewError(http.StatusBadRequest, "WL1002", fmt.Errorf("invalid wallet owner"))
	InvalidWalletBalanceOperation = error.NewError(http.StatusBadRequest, "WL1003", fmt.Errorf("invalid wallet balance operation"))
	NotEnoughBalance              = error.NewError(http.StatusBadRequest, "WL1004", fmt.Errorf("not enough balance"))
)
