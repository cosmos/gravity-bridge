package main

import (
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

func TestValidatorOut(t *testing.T) {
	withPristineE2EEnvironment(t, func(
		wd string,
		pool *dockertest.Pool,
		network *dockertest.Network,
	) {
		err := pool.RemoveContainerByName("gravity1")
		require.NoError(t, err, "error removing gravity1")

		buildAndRunTestRunner(t, wd, pool, network, "VALIDATOR_OUT")
	})
}
