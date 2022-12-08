package daemon

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

func Migrate(containerId string, predump bool, dockerdAddrSrc string, dockerdAddrDst string) error {

	// initiate connection with src dockerd
	cliSrc, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(dockerdAddrSrc))
	if err != nil {
		return err
	}
	// initiate connection with dst dockerd
	cliDst, err := client.NewClientWithOpts(client.FromEnv, client.WithHost(dockerdAddrDst))
	if err != nil {
		return err
	}

	if predump {
		return migratePredump(containerId, cliSrc, cliDst)
	} else {

	}

	return nil

}

// two-pass migrate, first time checkpoint leaves container running
// the second time freezes the container and migrate
func migratePredump(containerId string, cliSrc *client.Client, cliDst *client.Client) error {
	errs := make(chan error, 10)
	done := make(chan int, 10)

	// first checkpoint for source container
	go func() {
		err := cliSrc.CheckpointCreate(context.Background(),
			containerId,
			types.CheckpointCreateOptions{Predump: true, Exit: false, CheckpointID: "predumpCheckpointA"})
		if err != nil {
			errs <- err
		}
		done <- 0
	}()

	// find the docker image and configs
	inspectJson, err := cliSrc.ContainerInspect(context.Background(),
		containerId)
	if err != nil {
		return err
	}

	// create a container in the destination container
	go func() {
		_, dstErr := cliDst.ContainerCreate(context.Background(),
			inspectJson.Config,
			inspectJson.HostConfig,
			&network.NetworkingConfig{EndpointsConfig: inspectJson.NetworkSettings.Networks},
			&specs.Platform{
				OS:           "linux",
				Architecture: "amd64",
			},
			containerId)
		if dstErr != nil {
			errs <- dstErr
		}
		done <- 0
	}()

	// sync
	<-done
	<-done

	// second time checkpoint the source container
	srcErr := cliSrc.CheckpointCreate(context.Background(),
		containerId,
		types.CheckpointCreateOptions{
			Predump:            true,
			Exit:               true,
			ParentCheckpointID: "predumpCheckpointA",
			CheckpointID:       "predumpCheckpointB"})
	if srcErr != nil {
		return srcErr
	}

	// scp the checkpoints to the new server

}
