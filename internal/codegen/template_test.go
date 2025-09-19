package codegen

import (
	"strings"
	"testing"
)

func TestBind(t *testing.T) {
	structs := map[string]*tmplStruct{
		"RentalProposal": {
			Name: "RentalProposal",
			Fields: []*tmplField{
				{Name: "landlord", Type: "string"},
				{Name: "tenant", Type: "string"},
				{Name: "terms", Type: "string"},
			},
		},
		"RentalAgreement": {
			Name: "RentalAgreement",
			Fields: []*tmplField{
				{Name: "landlord", Type: "string"},
				{Name: "tenant", Type: "string"},
				{Name: "terms", Type: "string"},
			},
		},
	}

	result, err := Bind("main", structs)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// Check that the result contains expected content
	if !strings.Contains(result, "package main") {
		t.Error("Generated code does not contain correct package declaration")
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

	// Print the result for inspection
	t.Logf("Generated code:\n%s", result)
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
