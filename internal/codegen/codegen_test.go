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
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "20a17897a6664ecb8a4dd3e10b384c8cc41181d26ecbb446c2d65ae0928686c9"

type PARTY string
type TEXT string
type INT64 int64
type BOOL bool
type DECIMAL *big.Int
type NUMERIC *big.Int
type DATE time.Time
type TIMESTAMP time.Time
type UNIT struct{}
type LIST []string
type MAP map[string]interface{}
type OPTIONAL *interface{}

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

// RentalProposal is a Template type
type RentalProposal struct {
	Landlord PARTY ` + "`json:\"landlord\"`" + `
	Tenant   PARTY ` + "`json:\"tenant\"`" + `
	Terms    TEXT  ` + "`json:\"terms\"`" + `
}
`

	require.Equal(t, expectedCode, res, "generated code should match expected output")
}

func TestGetMainDalfAllTypes(t *testing.T) {
	srcPath := "../../test-data/archives/2.9.1/Test.dar"
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

	// Validate the full generated code from real DAML structures
	expectedMainCode := `package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "e2d906db3930143bfa53f43c7a69c218c8b499c03556485f312523090684ff34"

type PARTY string
type TEXT string
type INT64 int64
type BOOL bool
type DECIMAL *big.Int
type NUMERIC *big.Int
type DATE time.Time
type TIMESTAMP time.Time
type UNIT struct{}
type LIST []string
type MAP map[string]interface{}
type OPTIONAL *interface{}

// Address is a variant/union type
type Address struct {
	Us *USAddress ` + "`json:\"US,omitempty\"`" + `
	Uk *UKAddress ` + "`json:\"UK,omitempty\"`" + `
}

// MarshalJSON implements custom JSON marshaling for Address
func (v Address) MarshalJSON() ([]byte, error) {

	if v.Us != nil {
		return json.Marshal(map[string]interface{}{
			"tag":   "US",
			"value": v.Us,
		})
	}

	if v.Uk != nil {
		return json.Marshal(map[string]interface{}{
			"tag":   "UK",
			"value": v.Uk,
		})
	}

	return json.Marshal(map[string]interface{}{})
}

// UnmarshalJSON implements custom JSON unmarshaling for Address
func (v *Address) UnmarshalJSON(data []byte) error {
	var tagged struct {
		Tag   string          ` + "`json:\"tag\"`" + `
		Value json.RawMessage ` + "`json:\"value\"`" + `
	}

	if err := json.Unmarshal(data, &tagged); err != nil {
		return err
	}

	switch tagged.Tag {

	case "US":
		var value USAddress
		if err := json.Unmarshal(tagged.Value, &value); err != nil {
			return err
		}
		v.Us = &value

	case "UK":
		var value UKAddress
		if err := json.Unmarshal(tagged.Value, &value); err != nil {
			return err
		}
		v.Uk = &value

	default:
		return fmt.Errorf("unknown tag: %s", tagged.Tag)
	}

	return nil
}

// American is a Template type
type American struct {
	Person  PARTY     ` + "`json:\"person\"`" + `
	Address USAddress ` + "`json:\"address\"`" + `
}

// Briton is a Template type
type Briton struct {
	Person  PARTY     ` + "`json:\"person\"`" + `
	Address UKAddress ` + "`json:\"address\"`" + `
}

// OptionalFields is a Template type
type OptionalFields struct {
	Party  PARTY    ` + "`json:\"party\"`" + `
	AMaybe OPTIONAL ` + "`json:\"aMaybe\"`" + `
}

// OptionalFieldsCleanUp is a Record type
type OptionalFieldsCleanUp struct {
}

// Person is a Template type
type Person struct {
	Person  PARTY   ` + "`json:\"person\"`" + `
	Address Address ` + "`json:\"address\"`" + `
}

// SimpleFields is a Template type
type SimpleFields struct {
	Party     PARTY     ` + "`json:\"party\"`" + `
	ABool     BOOL      ` + "`json:\"aBool\"`" + `
	AInt      INT64     ` + "`json:\"aInt\"`" + `
	ADecimal  NUMERIC   ` + "`json:\"aDecimal\"`" + `
	AText     TEXT      ` + "`json:\"aText\"`" + `
	ADate     DATE      ` + "`json:\"aDate\"`" + `
	ADatetime TIMESTAMP ` + "`json:\"aDatetime\"`" + `
}

// SimpleFieldsCleanUp is a Record type
type SimpleFieldsCleanUp struct {
}

// UKAddress is a Record type
type UKAddress struct {
	Address  LIST     ` + "`json:\"address\"`" + `
	Locality OPTIONAL ` + "`json:\"locality\"`" + `
	City     TEXT     ` + "`json:\"city\"`" + `
	State    TEXT     ` + "`json:\"state\"`" + `
	Postcode TEXT     ` + "`json:\"postcode\"`" + `
}

// USAddress is a Record type
type USAddress struct {
	Address LIST  ` + "`json:\"address\"`" + `
	City    TEXT  ` + "`json:\"city\"`" + `
	State   TEXT  ` + "`json:\"state\"`" + `
	Zip     INT64 ` + "`json:\"zip\"`" + `
}
`

	require.Equal(t, expectedMainCode, res, "Generated main package code should match expected output")
}

func DeactivatedTestGetMainDalfV2_10_1(t *testing.T) {
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
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "28583fa5d41b7d36a1030152ab474228229d05ebcd398cc238d5b7a563293952"

type PARTY string
type TEXT string
type INT64 int64
type BOOL bool
type DECIMAL *big.Int
type NUMERIC *big.Int
type DATE time.Time
type TIMESTAMP time.Time
type UNIT struct{}
type LIST []string
type MAP map[string]interface{}
type OPTIONAL *interface{}

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

// RentalProposal is a Template type
type RentalProposal struct {
	Landlord PARTY ` + "`json:\"landlord\"`" + `
	Tenant   PARTY ` + "`json:\"tenant\"`" + `
	Terms    TEXT  ` + "`json:\"terms\"`" + `
}
`

	require.Equal(t, expectedCode, res, "generated code should match expected output")
}
