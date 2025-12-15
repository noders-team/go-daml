package testutil

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog/log"
)

const (
	damlSandboxVersion  = "3.5.0-snapshot.20251106.0"
	containerName       = "go-daml-test-canton"
	containerLabelKey   = "go-daml-test"
	containerLabelValue = "canton-sandbox"
)

var (
	once       sync.Once
	setupErr   error
	cl         *client.DamlBindingClient
	dockerPool *dockertest.Pool
	resDaml    *dockertest.Resource
	grpcAddr   string
	adminAddr  string
)

func Setup(ctx context.Context) error {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
		defer cancel()

		dockerPool, err := dockertest.NewPool("")
		if err != nil {
			log.Fatal().Err(err).Msg("Could not connect to docker")
		}

		if err := dockerPool.Client.Ping(); err != nil {
			log.Fatal().Err(err).Msg("Could not ping docker")
		}

		resDaml, grpcAddr, adminAddr = initDamlSandbox(ctx, dockerPool)

		builder := client.NewDamlClient("", grpcAddr).WithAdminAddress(adminAddr)
		if strings.HasSuffix(grpcAddr, ":443") {
			tlsConfig := client.TlsConfig{}
			builder = builder.WithTLSConfig(tlsConfig)
		}

		cl, err = builder.Build(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to build DAML client")
		}

		log.Info().Msg("Canton sandbox initialization complete, setting up test environment")

		testUser := "app-provider"
		users, err := cl.UserMng.ListUsers(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to list users")
		}

		userExists := false
		for _, u := range users {
			log.Info().Msgf("existing user: %s, primary party: %s", u.ID, u.PrimaryParty)
			if u.ID == testUser {
				userExists = true
				log.Info().Msgf("user %s already exists", testUser)
			}
		}

		if !userExists {
			log.Info().Msgf("creating user %s", testUser)

			log.Info().Msg("waiting for synchronizer connection before allocating party...")
			if err := waitForSynchronizerConnection(ctx, cl, 2*time.Minute); err != nil {
				log.Fatal().Err(err).Msg("synchronizer connection timeout")
			}

			partyDetails, err := cl.PartyMng.AllocateParty(ctx, "", nil, "")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to allocate party")
			}
			log.Info().Msgf("allocated party: %s", partyDetails.Party)

			user := &model.User{
				ID:           testUser,
				PrimaryParty: partyDetails.Party,
			}
			rights := []*model.Right{
				{Type: model.CanActAs{Party: partyDetails.Party}},
				{Type: model.CanReadAs{Party: partyDetails.Party}},
			}
			_, err = cl.UserMng.CreateUser(ctx, user, rights)
			if err != nil {
				log.Fatal().Err(err).Msgf("failed to create user %s", testUser)
			}
			log.Info().Msgf("created user %s with party %s", testUser, partyDetails.Party)
		}

		log.Info().Msg("Test environment ready, running tests")
	})
	return setupErr
}

func Teardown() {
	if dockerPool != nil {
		if resDaml != nil {
			log.Info().Str("container", containerName).Msg("Keeping Canton sandbox container for reuse")
		}
	}
}

func CleanupContainer() {
	if dockerPool != nil && resDaml != nil {
		if err := dockerPool.Purge(resDaml); err != nil {
			log.Error().Err(err).Msg("Could not purge Canton sandbox resource")
		} else {
			log.Info().Str("container", containerName).Msg("Removed Canton sandbox container")
		}
	}
}

func findExistingContainer(pool *dockertest.Pool) (*dockertest.Resource, error) {
	containers, err := pool.Client.ListContainers(docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"name": {containerName},
		},
	})
	if err != nil {
		return nil, err
	}

	var matchedContainer *docker.APIContainers
	for i := range containers {
		for _, name := range containers[i].Names {
			if name == "/"+containerName || name == containerName {
				matchedContainer = &containers[i]
				break
			}
		}
		if matchedContainer != nil {
			break
		}
	}

	if matchedContainer == nil {
		return nil, fmt.Errorf("no existing container found")
	}

	if matchedContainer.State != "running" {
		log.Warn().Str("state", matchedContainer.State).Msg("Found container but it's not running, removing it")
		err := pool.Client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    matchedContainer.ID,
			Force: true,
		})
		if err != nil {
			log.Warn().Err(err).Msg("Failed to remove stopped container")
		}
		return nil, fmt.Errorf("container not running, removed")
	}

	portMap := make(map[docker.Port][]docker.PortBinding)
	for _, p := range matchedContainer.Ports {
		port := docker.Port(fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		portMap[port] = []docker.PortBinding{
			{
				HostIP:   p.IP,
				HostPort: fmt.Sprintf("%d", p.PublicPort),
			},
		}
	}

	resource := &dockertest.Resource{
		Container: &docker.Container{
			ID:    matchedContainer.ID,
			Name:  matchedContainer.Names[0],
			State: docker.State{Running: true},
			NetworkSettings: &docker.NetworkSettings{
				Ports: portMap,
			},
		},
	}

	return resource, nil
}

func initDamlSandbox(ctx context.Context, dockerPool *dockertest.Pool) (*dockertest.Resource, string, string) {
	ledgerAPIPort := "6865"
	adminAPIPort := "6866"

	existingResource, err := findExistingContainer(dockerPool)
	if err == nil && existingResource != nil {
		log.Info().Str("container", containerName).Msg("Reusing existing Canton sandbox container")

		mappedLedgerPort := existingResource.GetPort(ledgerAPIPort + "/tcp")
		grpcAddr := fmt.Sprintf("127.0.0.1:%s", mappedLedgerPort)

		mappedAdminPort := existingResource.GetPort(adminAPIPort + "/tcp")
		adminAddr := fmt.Sprintf("127.0.0.1:%s", mappedAdminPort)

		log.Info().Msgf("Reusing Canton sandbox, Ledger API (gRPC) on %s", grpcAddr)

		if err := waitForPort(ctx, mappedLedgerPort, 30*time.Second); err != nil {
			log.Warn().Err(err).Msg("Existing container not responsive, creating new one")
		} else {
			return existingResource, grpcAddr, adminAddr
		}
	}

	cantonConfig := `canton {
  mediators {
    mediator1 {
      admin-api.port = 6869
    }
  }
  sequencers {
    sequencer1 {
      admin-api.port = 6868
      public-api.port = 6867
      sequencer {
        type = reference
        config.storage.type = memory
      }
      storage.type = memory
    }
  }
  participants {
    sandbox {
      storage.type = memory
      admin-api {
	    address = "0.0.0.0"
	  	port = 6866
	  }
      ledger-api {
        address = "0.0.0.0"
        port = 6865
        user-management-service.enabled = true
      }
    }
  }
}
`

	tmpDir, err := os.MkdirTemp("", "canton-config-*")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create temp dir for Canton config")
	}

	configPath := fmt.Sprintf("%s/canton.conf", tmpDir)
	if err := os.WriteFile(configPath, []byte(cantonConfig), 0o644); err != nil {
		log.Fatal().Err(err).Msg("Could not write Canton config")
	}

	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "digitalasset/daml-sdk",
		Tag:        damlSandboxVersion,
		Platform:   "linux/amd64",
		Name:       containerName,
		Cmd: []string{
			"daml",
			"sandbox",
			"-c", "/canton/canton.conf",
			"--debug",
		},
		ExposedPorts: []string{ledgerAPIPort + "/tcp", adminAPIPort + "/tcp"},
		Mounts:       []string{fmt.Sprintf("%s:/canton/canton.conf:ro", configPath)},
		Labels: map[string]string{
			containerLabelKey: containerLabelValue,
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = false
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "is already in use") {
			log.Warn().Err(err).Msg("Container already exists, attempting to reuse it")
			existingResource, findErr := findExistingContainer(dockerPool)
			if findErr == nil && existingResource != nil {
				mappedLedgerPort := existingResource.GetPort(ledgerAPIPort + "/tcp")
				grpcAddr := fmt.Sprintf("127.0.0.1:%s", mappedLedgerPort)
				mappedAdminPort := existingResource.GetPort(adminAPIPort + "/tcp")
				adminAddr := fmt.Sprintf("127.0.0.1:%s", mappedAdminPort)
				log.Info().Str("container", containerName).Msg("Successfully attached to existing container")
				return existingResource, grpcAddr, adminAddr
			}
			log.Fatal().Err(err).Msg("Container exists but could not attach to it")
		}
		log.Fatal().Err(err).Msg("Could not start DAML sandbox")
	}

	resource.Expire(600)

	mappedLedgerPort := resource.GetPort(ledgerAPIPort + "/tcp")
	grpcAddr := fmt.Sprintf("127.0.0.1:%s", mappedLedgerPort)

	log.Info().Msgf("DAML sandbox started, Ledger API (gRPC) on %s", grpcAddr)

	if err := waitForPort(ctx, mappedLedgerPort, 2*time.Minute); err != nil {
		log.Fatal().Err(err).Msgf("DAML sandbox Ledger API port %s not ready", mappedLedgerPort)
	}
	log.Info().Msgf("canton ledger API port %s is ready", adminAPIPort)

	adminAPIPort = resource.GetPort(adminAPIPort + "/tcp")
	adminAddr = fmt.Sprintf("127.0.0.1:%s", adminAPIPort)
	if err := waitForPort(ctx, adminAPIPort, 2*time.Minute); err != nil {
		log.Fatal().Err(err).Msgf("Canton admin API port %s not ready", adminAPIPort)
	}
	log.Info().Msgf("canton admin API port %s is ready", adminAPIPort)

	log.Info().Msg("port is open, waiting for Canton to fully initialize...")
	if err := waitForCantonReady(ctx, dockerPool, resource, 3*time.Minute); err != nil {
		log.Fatal().Err(err).Msg("Canton sandbox initialization timeout")
	}

	return resource, grpcAddr, adminAddr
}

func waitForPort(ctx context.Context, port string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	address := fmt.Sprintf("127.0.0.1:%s", port)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			conn.Close()
			log.Info().Msgf("Port %s is ready", port)
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for port %s", port)
}

func waitForCantonReady(ctx context.Context, pool *dockertest.Pool, resource *dockertest.Resource, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	readyMessage := "Canton sandbox is ready"

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		buf := &logBuffer{}
		err := pool.Client.Logs(docker.LogsOptions{
			Container:    resource.Container.ID,
			OutputStream: buf,
			Stdout:       true,
			Stderr:       true,
			Tail:         "100",
		})
		if err == nil {
			if strings.Contains(buf.String(), readyMessage) {
				log.Info().Msg("Canton sandbox is ready")
				return nil
			}
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("sandbox timeout: Canton sandbox did not become ready within %v", timeout)
}

type logBuffer struct {
	data []byte
}

func (b *logBuffer) Write(p []byte) (n int, err error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *logBuffer) String() string {
	return string(b.data)
}

func GetClient() *client.DamlBindingClient {
	return cl
}

func GetAdminAddr() string {
	return adminAddr
}

func GetGrpcAddr() string {
	return grpcAddr
}

func waitForSynchronizerConnection(ctx context.Context, cl *client.DamlBindingClient, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		resp, err := cl.StateService.GetConnectedSynchronizers(ctx, &model.GetConnectedSynchronizersRequest{})
		if err == nil && resp != nil && len(resp.ConnectedSynchronizers) > 0 {
			log.Info().
				Int("count", len(resp.ConnectedSynchronizers)).
				Str("first_id", resp.ConnectedSynchronizers[0].SynchronizerID).
				Msg("synchronizer connection established")
			return nil
		}

		if err != nil {
			log.Debug().Err(err).Msg("checking for synchronizer connection")
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("no synchronizer connection after %v", timeout)
}
