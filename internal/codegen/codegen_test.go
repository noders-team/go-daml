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

	res, err := Bind("main", pkg.Structs)
	require.NoError(t, err)
	require.NotEmpty(t, res)
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
}
