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

	participantID, err := cl.PartyMng.GetParticipantID(context.Background())
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("participant ID: %s", participantID))

	response, err := cl.PartyMng.ListKnownParties(context.Background(), "", 10, "")
	if err != nil {
		panic(err)
	}
	for _, d := range response.PartyDetails {
		println(fmt.Sprintf("received party details: %+v", d))
	}

	allocDetails, err := cl.PartyMng.GetParties(context.Background(), []string{"participant_admin"}, "")
	if err != nil {
		panic(err)
	}
	for _, ad := range allocDetails {
		println(fmt.Sprintf("party details: %+v", ad))
	}

	newParty, err := cl.PartyMng.AllocateParty(context.Background(), fmt.Sprintf("test_party_%d", os.Getpid()), map[string]string{
		"description": "New test party",
		"type":        "testing",
	}, "")
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("new party allocated: %+v", newParty))

	updatedParty, err := cl.PartyMng.UpdatePartyDetails(context.Background(), &model.PartyDetails{
		Party:   newParty.Party,
		IsLocal: true,
		LocalMetadata: map[string]string{
			"updated":     "true",
			"version":     "v2",
			"description": "Updated test party",
		},
	}, &model.UpdateMask{
		Paths: []string{"local_metadata"},
	})
	if err != nil {
		fmt.Printf("update party details error: %v\n", err)
	} else {
		println(fmt.Sprintf("updated party: %+v", updatedParty))
	}

	err = cl.PartyMng.UpdatePartyIdentityProviderID(context.Background(), "participant_admin", "", "new-provider-id")
	if err != nil {
		fmt.Printf("update identity provider error (expected if not supported): %v\n", err)
	}

	updatedResponse, err := cl.PartyMng.ListKnownParties(context.Background(), "", 10, "")
	if err != nil {
		panic(err)
	}
	for _, d := range updatedResponse.PartyDetails {
		println(fmt.Sprintf("updated party details: %+v", d))
	}
}
