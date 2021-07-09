package main

import (
	"testing"

	"github.com/ory/dockertest/v3"
)

func TestV2HappyPath(t *testing.T) {
	withPristineE2EEnvironment(t, func(
		wd string,
		pool *dockertest.Pool,
		network *dockertest.Network,
	) {
		buildAndRunTestRunner(t, wd, pool, network, "V2_HAPPY_PATH")
	})
}
