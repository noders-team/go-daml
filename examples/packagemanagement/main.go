package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
)

func main() {
	grpcAddress := os.Getenv("GRPC_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = "localhost:8080"
	}

	bearerToken := os.Getenv("BEARER_TOKEN")
	if bearerToken == "" {
		fmt.Println("Warning: BEARER_TOKEN environment variable not set")
	}

	tlsConfig := client.TlsConfig{}

	cl, err := client.NewDamlClient(bearerToken, grpcAddress).
		WithTLSConfig(tlsConfig).
		Build(context.Background())
	if err != nil {
		panic(err)
	}

	packages, err := cl.PackageMng.ListKnownPackages(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("known packages:")
	for _, pkg := range packages {
		println(fmt.Sprintf("  Package ID: %s", pkg.PackageID))
		println(fmt.Sprintf("  Name: %s", pkg.Name))
		println(fmt.Sprintf("  Version: %s", pkg.Version))
		println(fmt.Sprintf("  Size: %d bytes", pkg.PackageSize))
		if pkg.KnownSince != nil {
			println(fmt.Sprintf("  Known Since: %s", pkg.KnownSince.Format(time.RFC3339)))
		}
		println("")
	}

	darFilePath := "./examples/packagemanagement/test.dar"
	fmt.Printf("\ntesting DAR file upload from: %s\n", darFilePath)

	darContent, err := os.ReadFile(darFilePath)
	if err != nil {
		fmt.Printf("failed to read DAR file: %v\n", err)
	} else {
		fmt.Printf("dar file size: %d bytes\n", len(darContent))

		submissionID := fmt.Sprintf("validate-%d", time.Now().Unix())
		fmt.Printf("validating DAR file (submission ID: %s)...\n", submissionID)

		err = cl.PackageMng.ValidateDarFile(context.Background(), darContent, submissionID)
		if err != nil {
			fmt.Printf("dar validation failed: %v\n", err)
		} else {
			fmt.Println("dar validation successful!")

			uploadSubmissionID := fmt.Sprintf("upload-%d", time.Now().Unix())
			fmt.Printf("uploading DAR file (submission ID: %s)...\n", uploadSubmissionID)

			err = cl.PackageMng.UploadDarFile(context.Background(), darContent, uploadSubmissionID)
			if err != nil {
				fmt.Printf("dar upload failed: %v\n", err)
			} else {
				fmt.Println("dar upload successful!")

				updatedPackages, err := cl.PackageMng.ListKnownPackages(context.Background())
				if err != nil {
					fmt.Printf("failed to list packages after upload: %v\n", err)
				} else {
					fmt.Printf("total packages after upload: %d\n", len(updatedPackages))
				}
			}
		}
	}
}
