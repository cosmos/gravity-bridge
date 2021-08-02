package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func fileCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func writeFile(t *testing.T, path string, body []byte) {
	t.Helper()

	_, err := os.Create(path)
	require.NoError(t, err)

	err = ioutil.WriteFile(path, body, 0644)
	require.NoError(t, err)
}

func buildAndRunTestRunner(t *testing.T,
	wd string,
	pool *dockertest.Pool,
	network *dockertest.Network,
	testType string,
) {
	t.Helper()

	// bring up the test runner
	t.Log("building and deploying test runner")
	testRunner, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       "test_runner",
			Repository: "test-runner",
			Tag:        "prebuilt",
			NetworkID:  network.Network.ID,
			PortBindings: map[docker.Port][]docker.PortBinding{
				"8545/tcp": {{HostIP: "", HostPort: "8545"}},
			},
			Mounts: []string{
				fmt.Sprintf("%s/testdata:/testdata", wd),
			},
			Env: []string{
				"RUST_BACKTRACE=full",
				"RUST_LOG=INFO",
				fmt.Sprintf("TEST_TYPE=%s", testType),
			},
		}, func(config *docker.HostConfig) {})

	require.NoError(t, err, "error bringing up test runner")
	t.Logf("deployed test runner at %s", testRunner.Container.ID)

	container := testRunner.Container
	for container.State.Running {
		time.Sleep(10 * time.Second)
		container, err = pool.Client.InspectContainer(container.ID)
		require.NoError(t, err, "error inspecting test runner")
	}
	require.Equal(t, 0, container.State.ExitCode, "container exited with error")
}
