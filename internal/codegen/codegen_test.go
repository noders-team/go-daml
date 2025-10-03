package codegen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMainDalf(t *testing.T) {
	srcPath := "../../test-data/test.dar"
	output := "../../test-data/test_unzipped"
	defer os.RemoveAll(output)

	genOutput, err := UnzipDar(srcPath, &output)
	require.NoError(t, err)

	manifest, err := GetManifest(genOutput)
	require.NoError(t, err)
	require.Equal(t, "rental-0.1.0-20a17897a6664ecb8a4dd3e10b384c8cc41181d26ecbb446c2d65ae0928686c9/rental-0.1.0-20a17897a6664ecb8a4dd3e10b384c8cc41181d26ecbb446c2d65ae0928686c9.dalf", manifest.MainDalf)
	require.NotNil(t, manifest)
	require.Equal(t, "1.0", manifest.Version)
	require.Equal(t, "damlc", manifest.CreatedBy)
	require.Equal(t, "rental-0.1.0", manifest.Name)
	require.Equal(t, "1.18.1", manifest.SdkVersion)
	require.Equal(t, "daml-lf", manifest.Format)
	require.Equal(t, "non-encrypted", manifest.Encryption)
	require.Len(t, manifest.Dalfs, 25)

	dalfFullPath := filepath.Join(genOutput, manifest.MainDalf)
	dalfContent, err := os.ReadFile(dalfFullPath)
	require.NoError(t, err)
	require.NotNil(t, dalfContent)

	pkg, err := GetAST(dalfContent, manifest)
	require.Nil(t, err)
	require.NotEmpty(t, pkg.Structs)

	pkg1, exists := pkg.Structs["RentalAgreement"]
	require.True(t, exists)
	require.Len(t, pkg1.Fields, 3)
	require.Equal(t, pkg1.Name, "RentalAgreement")
	require.Equal(t, pkg1.Fields[0].Name, "landlord")
	require.Equal(t, pkg1.Fields[1].Name, "tenant")
	require.Equal(t, pkg1.Fields[2].Name, "terms")

	pkg2, exists := pkg.Structs["Accept"]
	require.True(t, exists)
	require.Len(t, pkg2.Fields, 2)
	require.Equal(t, pkg2.Name, "Accept")
	require.Equal(t, pkg2.Fields[0].Name, "foo")
	require.Equal(t, pkg2.Fields[1].Name, "bar")

	pkg3, exists := pkg.Structs["RentalProposal"]
	require.True(t, exists)
	require.Len(t, pkg3.Fields, 3)
	require.Equal(t, pkg3.Name, "RentalProposal")
	require.Equal(t, pkg3.Fields[0].Name, "landlord")
	require.Equal(t, pkg3.Fields[1].Name, "tenant")
	require.Equal(t, pkg3.Fields[2].Name, "terms")

	res, err := Bind("main", pkg.PackageID, pkg.Structs)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	// Validate the full generated code
	expectedCode := `package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "20a17897a6664ecb8a4dd3e10b384c8cc41181d26ecbb446c2d65ae0928686c9"

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	return map[string]interface{}{
		"args": args,
	}
}

// Accept is a Record type
type Accept struct {
	Foo TEXT  ` + "`json:\"foo\"`" + `
	Bar INT64 ` + "`json:\"bar\"`" + `
}

// RentalAgreement is a Template type
type RentalAgreement struct {
	Landlord PARTY ` + "`json:\"landlord\"`" + `
	Tenant   PARTY ` + "`json:\"tenant\"`" + `
	Terms    TEXT  ` + "`json:\"terms\"`" + `
}

// GetTemplateID returns the template ID for this template
func (t RentalAgreement) GetTemplateID() string {
	return fmt.Sprintf("%s:%s:%s", PackageID, "Rental", "RentalAgreement")
}

// CreateCommand returns a CreateCommand for this template
func (t RentalAgreement) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["landlord"] = map[string]interface{}{"_type": "party", "value": string(t.Landlord)}

	args["tenant"] = map[string]interface{}{"_type": "party", "value": string(t.Tenant)}

	args["terms"] = string(t.Terms)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// Choice methods for RentalAgreement

// Archive exercises the Archive choice on this RentalAgreement contract
func (t RentalAgreement) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "Rental", "RentalAgreement"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// RentalProposal is a Template type
type RentalProposal struct {
	Landlord PARTY ` + "`json:\"landlord\"`" + `
	Tenant   PARTY ` + "`json:\"tenant\"`" + `
	Terms    TEXT  ` + "`json:\"terms\"`" + `
}

// GetTemplateID returns the template ID for this template
func (t RentalProposal) GetTemplateID() string {
	return fmt.Sprintf("%s:%s:%s", PackageID, "Rental", "RentalProposal")
}

// CreateCommand returns a CreateCommand for this template
func (t RentalProposal) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["landlord"] = map[string]interface{}{"_type": "party", "value": string(t.Landlord)}

	args["tenant"] = map[string]interface{}{"_type": "party", "value": string(t.Tenant)}

	args["terms"] = string(t.Terms)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// Choice methods for RentalProposal

// Archive exercises the Archive choice on this RentalProposal contract
func (t RentalProposal) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "Rental", "RentalProposal"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// Accept exercises the Accept choice on this RentalProposal contract
func (t RentalProposal) Accept(contractID string, args Accept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", PackageID, "Rental", "RentalProposal"),
		ContractID: contractID,
		Choice:     "Accept",
		Arguments:  argsToMap(args),
	}
}
`

	require.Equal(t, expectedCode, res, "generated code should match expected output")
}

func TestGetMainDalfAllTypes(t *testing.T) {
	srcPath := "../../test-data/test_2_9_1.dar"
	output := "../../test-data/test_unzipped"
	defer os.RemoveAll(output)

	genOutput, err := UnzipDar(srcPath, &output)
	require.NoError(t, err)
	defer os.RemoveAll(genOutput)

	manifest, err := GetManifest(genOutput)
	require.NoError(t, err)
	require.Equal(t, "Test-1.0.0-e2d906db3930143bfa53f43c7a69c218c8b499c03556485f312523090684ff34/Test-1.0.0-e2d906db3930143bfa53f43c7a69c218c8b499c03556485f312523090684ff34.dalf", manifest.MainDalf)
	require.NotNil(t, manifest)
	require.Equal(t, "1.0", manifest.Version)
	require.Equal(t, "damlc", manifest.CreatedBy)
	require.Equal(t, "Test-1.0.0", manifest.Name)
	require.Equal(t, "2.9.1", manifest.SdkVersion)
	require.Equal(t, "daml-lf", manifest.Format)
	require.Equal(t, "non-encrypted", manifest.Encryption)
	require.Len(t, manifest.Dalfs, 29)

	dalfFullPath := filepath.Join(genOutput, manifest.MainDalf)
	dalfContent, err := os.ReadFile(dalfFullPath)
	require.NoError(t, err)
	require.NotNil(t, dalfContent)

	pkg, err := GetAST(dalfContent, manifest)
	require.Nil(t, err)
	require.NotEmpty(t, pkg.Structs)

	// Test Address struct (variant/union type)
	addressStruct, exists := pkg.Structs["Address"]
	require.True(t, exists)
	require.Len(t, addressStruct.Fields, 2)
	require.Equal(t, addressStruct.Name, "Address")
	require.Equal(t, addressStruct.Fields[0].Name, "US")
	require.Equal(t, addressStruct.Fields[0].Type, "USAddress")
	require.Equal(t, addressStruct.Fields[1].Name, "UK")
	require.Equal(t, addressStruct.Fields[1].Type, "UKAddress")

	// Test USAddress struct
	usAddressStruct, exists := pkg.Structs["USAddress"]
	require.True(t, exists)
	require.Len(t, usAddressStruct.Fields, 4)
	require.Equal(t, usAddressStruct.Name, "USAddress")
	require.Equal(t, usAddressStruct.Fields[0].Name, "address")
	require.Equal(t, usAddressStruct.Fields[1].Name, "city")
	require.Equal(t, usAddressStruct.Fields[2].Name, "state")
	require.Equal(t, usAddressStruct.Fields[3].Name, "zip")

	// Test UKAddress struct
	ukAddressStruct, exists := pkg.Structs["UKAddress"]
	require.True(t, exists)
	require.Len(t, ukAddressStruct.Fields, 5)
	require.Equal(t, ukAddressStruct.Name, "UKAddress")
	require.Equal(t, ukAddressStruct.Fields[0].Name, "address")
	require.Equal(t, ukAddressStruct.Fields[1].Name, "locality")
	require.Equal(t, ukAddressStruct.Fields[2].Name, "city")
	require.Equal(t, ukAddressStruct.Fields[3].Name, "state")
	require.Equal(t, ukAddressStruct.Fields[4].Name, "postcode")

	// Test Person struct (uses Address)
	personStruct, exists := pkg.Structs["Person"]
	require.True(t, exists)
	require.Len(t, personStruct.Fields, 2)
	require.Equal(t, personStruct.Name, "Person")
	require.Equal(t, personStruct.Fields[0].Name, "person")
	require.Equal(t, personStruct.Fields[1].Name, "address")
	require.Equal(t, personStruct.Fields[1].Type, "Address")

	// Test American struct (uses USAddress)
	americanStruct, exists := pkg.Structs["American"]
	require.True(t, exists)
	require.Len(t, americanStruct.Fields, 2)
	require.Equal(t, americanStruct.Name, "American")
	require.Equal(t, americanStruct.Fields[0].Name, "person")
	require.Equal(t, americanStruct.Fields[1].Name, "address")
	require.Equal(t, americanStruct.Fields[1].Type, "USAddress")

	// Test Briton struct (uses UKAddress)
	britonStruct, exists := pkg.Structs["Briton"]
	require.True(t, exists)
	require.Len(t, britonStruct.Fields, 2)
	require.Equal(t, britonStruct.Name, "Briton")
	require.Equal(t, britonStruct.Fields[0].Name, "person")
	require.Equal(t, britonStruct.Fields[1].Name, "address")
	require.Equal(t, britonStruct.Fields[1].Type, "UKAddress")

	// Test SimpleFields struct (various primitive types)
	simpleFieldsStruct, exists := pkg.Structs["SimpleFields"]
	require.True(t, exists)
	require.Len(t, simpleFieldsStruct.Fields, 7)
	require.Equal(t, simpleFieldsStruct.Name, "SimpleFields")
	require.Equal(t, simpleFieldsStruct.Fields[0].Name, "party")
	require.Equal(t, simpleFieldsStruct.Fields[1].Name, "aBool")
	require.Equal(t, simpleFieldsStruct.Fields[2].Name, "aInt")
	require.Equal(t, simpleFieldsStruct.Fields[3].Name, "aDecimal")
	require.Equal(t, simpleFieldsStruct.Fields[4].Name, "aText")
	require.Equal(t, simpleFieldsStruct.Fields[5].Name, "aDate")
	require.Equal(t, simpleFieldsStruct.Fields[6].Name, "aDatetime")

	// Test OptionalFields struct
	optionalFieldsStruct, exists := pkg.Structs["OptionalFields"]
	require.True(t, exists)
	require.Len(t, optionalFieldsStruct.Fields, 2)
	require.Equal(t, optionalFieldsStruct.Name, "OptionalFields")
	require.Equal(t, optionalFieldsStruct.Fields[0].Name, "party")
	require.Equal(t, optionalFieldsStruct.Fields[1].Name, "aMaybe")

	// Test that Address struct is identified as variant
	require.Equal(t, "Variant", addressStruct.RawType, "Address should be identified as variant type")
	require.True(t, addressStruct.Fields[0].IsOptional, "US field should be optional")
	require.True(t, addressStruct.Fields[1].IsOptional, "UK field should be optional")

	// Test that non-variant structs have correct RawType
	require.Equal(t, "Record", usAddressStruct.RawType, "USAddress should be Record type")
	require.Equal(t, "Record", ukAddressStruct.RawType, "UKAddress should be Record type")
	// Note: Some structs might be templates in the new template-first approach
	if personStruct.RawType != "Record" && personStruct.RawType != "Template" {
		require.Fail(t, "Person should be either Record or Template type, got: %s", personStruct.RawType)
	}
	if americanStruct.RawType != "Record" && americanStruct.RawType != "Template" {
		require.Fail(t, "American should be either Record or Template type, got: %s", americanStruct.RawType)
	}
	if britonStruct.RawType != "Record" && britonStruct.RawType != "Template" {
		require.Fail(t, "Briton should be either Record or Template type, got: %s", britonStruct.RawType)
	}
	if simpleFieldsStruct.RawType != "Record" && simpleFieldsStruct.RawType != "Template" {
		require.Fail(t, "SimpleFields should be either Record or Template type, got: %s", simpleFieldsStruct.RawType)
	}
	if optionalFieldsStruct.RawType != "Record" && optionalFieldsStruct.RawType != "Template" {
		require.Fail(t, "OptionalFields should be either Record or Template type, got: %s", optionalFieldsStruct.RawType)
	}

	res, err := Bind("main", pkg.PackageID, pkg.Structs)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	testData2_9_1 := "../../test-data/test_1_0_0.go_gen"
	expectedMainCode, err := os.ReadFile(testData2_9_1)
	require.NoError(t, err)

	// Validate the full generated code from real DAML structures
	require.Equal(t, string(expectedMainCode), res, "Generated main package code should match expected output")
}

func TestGetMainDalfV3(t *testing.T) {
	srcPath := "../../test-data/all-kinds-of-1.0.0_lf.dar"
	output := "../../test-data/test_unzipped"
	defer os.RemoveAll(output)

	genOutput, err := UnzipDar(srcPath, &output)
	require.NoError(t, err)

	manifest, err := GetManifest(genOutput)
	require.NoError(t, err)
	require.Equal(t, "all-kinds-of-1.0.0-6d7e83e81a0a7960eec37340f5b11e7a61606bd9161f413684bc345c3f387948/all-kinds-of-1.0.0-6d7e83e81a0a7960eec37340f5b11e7a61606bd9161f413684bc345c3f387948.dalf", manifest.MainDalf)
	require.NotNil(t, manifest)
	require.Equal(t, "1.0", manifest.Version)
	require.Equal(t, "damlc", manifest.CreatedBy)
	require.Equal(t, "all-kinds-of-1.0.0", manifest.Name)
	require.Equal(t, "3.3.0-snapshot.20250417.0", manifest.SdkVersion)
	require.Equal(t, "daml-lf", manifest.Format)
	require.Equal(t, "non-encrypted", manifest.Encryption)
	require.Len(t, manifest.Dalfs, 30)

	dalfFullPath := filepath.Join(genOutput, manifest.MainDalf)
	dalfContent, err := os.ReadFile(dalfFullPath)
	require.NoError(t, err)
	require.NotNil(t, dalfContent)

	pkg, err := GetAST(dalfContent, manifest)
	require.Nil(t, err)
	require.NotEmpty(t, pkg.Structs)

	// Test MappyContract template
	pkg1, exists := pkg.Structs["MappyContract"]
	require.True(t, exists)
	require.Equal(t, pkg1.Name, "MappyContract")
	require.Equal(t, "Template", pkg1.RawType)
	require.Len(t, pkg1.Fields, 2)
	require.Equal(t, pkg1.Fields[0].Name, "operator")
	require.Equal(t, pkg1.Fields[1].Name, "value")

	// Test OneOfEverything template
	pkg2, exists := pkg.Structs["OneOfEverything"]
	require.True(t, exists)
	require.Equal(t, pkg2.Name, "OneOfEverything")
	require.Equal(t, "Template", pkg2.RawType)
	require.Len(t, pkg2.Fields, 16) // Based on the generated output
	require.Equal(t, pkg2.Fields[0].Name, "operator")
	require.Equal(t, pkg2.Fields[1].Name, "someBoolean")
	require.Equal(t, pkg2.Fields[2].Name, "someInteger")

	// Test Accept struct
	pkg3, exists := pkg.Structs["Accept"]
	require.True(t, exists)
	require.Equal(t, pkg3.Name, "Accept")
	require.Equal(t, "Record", pkg3.RawType)

	// Test Color enum
	colorStruct, exists := pkg.Structs["Color"]
	require.True(t, exists)
	require.Equal(t, "Enum", colorStruct.RawType)
	require.Len(t, colorStruct.Fields, 3)
	require.Equal(t, colorStruct.Fields[0].Name, "Red")
	require.Equal(t, colorStruct.Fields[1].Name, "Green")
	require.Equal(t, colorStruct.Fields[2].Name, "Blue")

	res, err := Bind("main", pkg.PackageID, pkg.Structs)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	testRes := "../../test-data/all_kinds_of_1_0_0.go_gen"
	expectedCode, err := os.ReadFile(testRes)
	require.NoError(t, err)

	require.Equal(t, string(expectedCode), res, "generated code should match expected output")
}
