package codegen

import (
	"strings"
	"testing"

	"github.com/noders-team/go-daml/internal/codegen/model"
)

func TestBind(t *testing.T) {
	structs := map[string]*model.TmplStruct{
		"RentalProposal": {
			Name:    "RentalProposal",
			RawType: "Record",
			Fields: []*model.TmplField{
				{Name: "landlord", Type: "string"},
				{Name: "tenant", Type: "string"},
				{Name: "terms", Type: "string"},
			},
		},
		"RentalAgreement": {
			Name:    "RentalAgreement",
			RawType: "Record",
			Fields: []*model.TmplField{
				{Name: "landlord", Type: "string"},
				{Name: "tenant", Type: "string"},
				{Name: "terms", Type: "string"},
			},
		},
	}

	result, err := Bind("main", "test-package-id", "2.0.0", structs)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// Check that the result contains expected content
	if !strings.Contains(result, "package main") {
		t.Error("Generated code does not contain correct package declaration")
	}

	if !strings.Contains(result, `const PackageID = "test-package-id"`) {
		t.Error("Generated code does not contain PackageID constant")
	}

	if !strings.Contains(result, "type RentalProposal struct") {
		t.Error("Generated code does not contain RentalProposal struct")
	}

	if !strings.Contains(result, "type RentalAgreement struct") {
		t.Error("Generated code does not contain RentalAgreement struct")
	}

	if !strings.Contains(result, "Landlord string") {
		t.Error("Generated code does not contain capitalized field names")
	}

	if !strings.Contains(result, `json:"landlord"`) {
		t.Error("Generated code does not contain JSON tags with original field names")
	}

	if !strings.Contains(result, `const SDKVersion = "2.0.0"`) {
		t.Error("Generated code does not contain SDKVersion constant")
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"landlord", "Landlord"},
		{"rental_proposal", "RentalProposal"},
		{"TENANT", "Tenant"},
		{"", ""},
		{"a", "A"},
	}

	for _, test := range tests {
		result := capitalize(test.input)
		if result != test.expected {
			t.Errorf("capitalize(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestDecapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Landlord", "landlord"},
		{"RentalProposal", "rentalProposal"},
		{"TENANT", "tenant"},
		{"", ""},
		{"A", "a"},
	}

	for _, test := range tests {
		result := decapitalize(test.input)
		if result != test.expected {
			t.Errorf("decapitalize(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
