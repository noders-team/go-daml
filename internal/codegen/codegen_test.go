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

	_, err := UnzipDar(srcPath, &output)
	require.NoError(t, err)

	manifest, err := GetManifest(output)
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

	dalfFullPath := filepath.Join(output, manifest.MainDalf)
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

	//res, err := Bind("main", pkg.Structs)
	//require.NoError(t, err)
	//require.NotEmpty(t, res)
}

func TestGetMainDalfV2(t *testing.T) {
	srcPath := "../../test-data/archives/2.9.1/Test.dar"
	output := "../../test-data/test_unzipped"
	defer os.RemoveAll(output)

	resDir, err := UnzipDar(srcPath, &output)
	require.NoError(t, err)
	defer os.RemoveAll(resDir)

	manifest, err := GetManifest(output)
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

	dalfFullPath := filepath.Join(output, manifest.MainDalf)
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
}
