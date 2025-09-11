package main

import (
	"context"
	"fmt"
	"os"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
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

	users, err := cl.UserMng.ListUsers(context.Background())
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		println(fmt.Sprintf("received user details: %+v ", u))
	}

	user, err := cl.UserMng.GetUser(context.Background(), "participant_admin")
	if err != nil {
		panic(err)
	}

	println(fmt.Sprintf("single user details: %+v ", user))

	userRights, err := cl.UserMng.ListUserRights(context.Background(), user.ID)
	if err != nil {
		panic(err)
	}
	for _, r := range userRights {
		println(fmt.Sprintf("user rights: %+v ", r))
	}

	newRights := make([]*model.Right, 0)
	newRights = append(newRights, &model.Right{Type: model.CanReadAs{Party: "app_provider_localnet-localparty-1"}})

	updatedRights, err := cl.UserMng.GrantUserRights(context.Background(), user.ID, newRights)
	if err != nil {
		panic(err)
	}
	for _, r := range updatedRights {
		println(fmt.Sprintf("user rights: %+v ", r))
	}

	updatedRights, err = cl.UserMng.RevokeUserRights(context.Background(), user.ID, newRights)
	if err != nil {
		panic(err)
	}
	for _, r := range updatedRights {
		println(fmt.Sprintf("user rights: %+v ", r))
	}
}
