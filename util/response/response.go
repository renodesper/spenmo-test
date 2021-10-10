package response

import (
	"time"

	"github.com/google/uuid"
	e "gitlab.com/renodesper/spenmo-test/util/error"
)

type (
	// ErrorResponse ...
	ErrorResponse struct {
		Errors []e.Error `json:"errors"`
		Meta   e.Meta    `json:"meta"`
	}

	// SuccessResponse ...
	SuccessResponse struct {
		Data interface{} `json:"data"`
		Meta e.Meta      `json:"meta"`
	}
)

// PopulateMeta will return current timestamp and requestId as meta
func PopulateMeta(requestID string) e.Meta {
	if requestID == "" {
		requestID = uuid.New().String()
	}

	return e.Meta{
		"timestamp": time.Now().UnixNano(),
		"requestId": requestID,
	}
}
