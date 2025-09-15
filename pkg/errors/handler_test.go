package errors

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAsDamlError(t *testing.T) {
	tests := []struct {
		name          string
		inputError    error
		expectedError *DamlError
		shouldFail    bool
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
			shouldFail: false,
		},
		{
			name:          "nil error",
			inputError:    nil,
			expectedError: nil,
			shouldFail:    true,
		},
		{
			name:          "non-gRPC error",
			inputError:    status.Error(codes.Internal, "regular error message without DAML format"),
			expectedError: nil,
			shouldFail:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AsDamlError(tt.inputError)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
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
