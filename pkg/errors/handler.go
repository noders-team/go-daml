package errors

import (
	"regexp"
	"strconv"

	"google.golang.org/grpc/status"
)

const genericErr = "DAML_GENERIC_ERROR_CODE"

type DamlError struct {
	ErrorCode     string
	CategoryID    int
	CorrelationID interface{}
	Message       string
}

func AsDamlError(err error) *DamlError {
	if err == nil {
		return &DamlError{
			ErrorCode:  genericErr,
			CategoryID: -1,
			Message:    err.Error(),
		}
	}

	grpcStatus, ok := status.FromError(err)
	if !ok {
		return &DamlError{
			ErrorCode:  genericErr,
			CategoryID: -2,
			Message:    err.Error(),
		}
	}

	damlErrorRegex := regexp.MustCompile(`^([A-Z_]+)\((\d+),([^)]+)\):\s*(.*)$`)
	message := grpcStatus.Message()

	matches := damlErrorRegex.FindStringSubmatch(message)
	if len(matches) == 5 {
		categoryID, err := strconv.Atoi(matches[2])
		if err != nil {
			return &DamlError{
				ErrorCode:  genericErr,
				CategoryID: -3,
				Message:    err.Error(),
			}
		}

		return &DamlError{
			ErrorCode:     matches[1],
			CategoryID:    categoryID,
			CorrelationID: matches[3],
			Message:       matches[4],
		}
	}

	return &DamlError{
		ErrorCode:  genericErr,
		CategoryID: -5,
		Message:    err.Error(),
	}
}
