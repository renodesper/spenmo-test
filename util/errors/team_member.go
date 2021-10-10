package errors

import (
	"fmt"
	"net/http"

	"gitlab.com/renodesper/spenmo-test/util/error"
)

var (
	MissingTeamID = error.NewError(http.StatusBadRequest, "TM1000", fmt.Errorf("missing teamId"))
	MissingUserID = error.NewError(http.StatusBadRequest, "TM1001", fmt.Errorf("missing userId"))
)
