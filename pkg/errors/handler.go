package errors

import (
	"fmt"
	"regexp"
	"strconv"

	"google.golang.org/grpc/status"
)

type DamlError struct {
	ErrorCode     string
	CategoryID    int
	CorrelationID interface{}
	Message       string
}

func AsDamlError(err error) (*DamlError, error) {
	if err == nil {
		return nil, fmt.Errorf("no error")
	}

	grpcStatus, ok := status.FromError(err)
	if !ok {
		return nil, err
	}

	damlErrorRegex := regexp.MustCompile(`^([A-Z_]+)\((\d+),([^)]+)\):\s*(.*)$`)
	message := grpcStatus.Message()

	matches := damlErrorRegex.FindStringSubmatch(message)
	if len(matches) == 5 {
		categoryID, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %v", err)
		}

		return &DamlError{
			ErrorCode:     matches[1],
			CategoryID:    categoryID,
			CorrelationID: matches[3],
			Message:       matches[4],
		}, nil
	}

	return nil, fmt.Errorf("unable to parse DAML error: %s", message)
}
