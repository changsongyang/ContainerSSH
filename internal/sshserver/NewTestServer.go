package sshserver

import (
	"context"
	"fmt"
	"testing"

    config2 "go.containerssh.io/containerssh/config"
    "go.containerssh.io/containerssh/internal/structutils"
    "go.containerssh.io/containerssh/internal/test"
    "go.containerssh.io/containerssh/log"
    "go.containerssh.io/containerssh/service"
)

// NewTestServer is a simplified API to start and stop a test server.
func NewTestServer(t *testing.T, handler Handler, logger log.Logger, config *config2.SSHConfig) TestServer {
	if config == nil {
		config = &config2.SSHConfig{}
		structutils.Defaults(config)
	}

	port := test.GetNextPort(t, "SSH")
	config.Listen = fmt.Sprintf("127.0.0.1:%d", port)
	if err := config.GenerateHostKey(); err != nil {
		panic(err)
	}
	svc, err := New(*config, handler, logger)
	if err != nil {
		panic(err)
	}
	lifecycle := service.NewLifecycle(svc)
	started := make(chan struct{})
	lifecycle.OnRunning(
		func(s service.Service, l service.Lifecycle) {
			started <- struct{}{}
		})

	t.Cleanup(func() {
		lifecycle.Stop(context.Background())
		_ = lifecycle.Wait()
	})

	return &testServerImpl{
		config:    *config,
		lifecycle: lifecycle,
		started:   started,
	}
}
