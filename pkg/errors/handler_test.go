package errors

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAsDamlError(t *testing.T) {
	tests := []struct {
		name          string
		inputError    error
		expectedError *DamlError
	}{
		{
			name:       "valid DAML error from gRPC",
			inputError: status.Error(codes.NotFound, "PACKAGE_NAMES_NOT_FOUND(11,aa8b050d): The following package names do not match upgradable packages uploaded on this participant: [DamlScript]."),
			expectedError: &DamlError{
				ErrorCode:     "PACKAGE_NAMES_NOT_FOUND",
				CategoryID:    11,
				CorrelationID: "aa8b050d",
				Message:       "The following package names do not match upgradable packages uploaded on this participant: [DamlScript].",
			},
		},
		{
			name:       "non-gRPC error",
			inputError: errors.New("regular error message"),
			expectedError: &DamlError{
				ErrorCode:  genericErr,
				CategoryID: -2,
				Message:    "regular error message",
			},
		},
		{
			name:       "gRPC error without DAML format",
			inputError: status.Error(codes.Internal, "regular error message without DAML format"),
			expectedError: &DamlError{
				ErrorCode:  genericErr,
				CategoryID: -5,
				Message:    "rpc error: code = Internal desc = regular error message without DAML format",
			},
		},
		{
			name:       "gRPC error without DAML format (no regex match)",
			inputError: status.Error(codes.NotFound, "INVALID_ERROR(invalid_id,aa8b050d): Test message"),
			expectedError: &DamlError{
				ErrorCode:  genericErr,
				CategoryID: -5,
				Message:    "rpc error: code = NotFound desc = INVALID_ERROR(invalid_id,aa8b050d): Test message",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AsDamlError(tt.inputError)

			if result == nil {
				t.Fatal("AsDamlError should never return nil")
			}

			if result.ErrorCode != tt.expectedError.ErrorCode {
				t.Errorf("ErrorCode mismatch: got %s, want %s", result.ErrorCode, tt.expectedError.ErrorCode)
			}

			if result.CategoryID != tt.expectedError.CategoryID {
				t.Errorf("CategoryID mismatch: got %d, want %d", result.CategoryID, tt.expectedError.CategoryID)
			}

			if result.CorrelationID != tt.expectedError.CorrelationID {
				t.Errorf("CorrelationID mismatch: got %v, want %v", result.CorrelationID, tt.expectedError.CorrelationID)
			}

			if result.Message != tt.expectedError.Message {
				t.Errorf("Message mismatch: got %s, want %s", result.Message, tt.expectedError.Message)
			}
		})
	}
}

func TestAsDamlErrorPanicWithNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("AsDamlError should panic when passed nil error")
		}
	}()

	AsDamlError(nil)
}
