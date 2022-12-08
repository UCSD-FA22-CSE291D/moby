package daemon

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"os/exec"
	"strings"
)

func (daemon *Daemon) Migrate(containerId string, options types.MigrateOptions) error {

	predump := options.Predump
	dockerdAddrSrc := options.SrcDockerdAddr
	dockerdAddrDst := options.DstDockerdAddr

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

	srcAddr := strings.Split(dockerdAddrSrc, ":")[0]
	dstAddr := strings.Split(dockerdAddrDst, ":")[0]

	if predump {
		return daemon.migratePredump(containerId, cliSrc, cliDst, srcAddr, dstAddr)
	} else {
		return daemon.migrate(containerId, cliSrc, cliDst, srcAddr, dstAddr)
	}

}

// two-pass migrate, first time checkpoint leaves container running
// the second time freezes the container and migrate
func (daemon *Daemon) migratePredump(containerId string, cliSrc *client.Client, cliDst *client.Client, srcAddr string, dstAddr string) error {
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
	createResp, dstErr := cliDst.ContainerCreate(context.Background(),
		inspectJson.Config,
		inspectJson.HostConfig,
		&network.NetworkingConfig{EndpointsConfig: inspectJson.NetworkSettings.Networks},
		&specs.Platform{
			OS:           "linux",
			Architecture: "amd64",
		},
		containerId)
	if dstErr != nil {
		return dstErr
	}

	dstContainerId := createResp.ID

	// sync
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
	cmd := exec.Command("scp %s:/var/lib/docker/containers/%s/checkpoints/predumpCheckpointB "+
		"%s:/var/lib/docker/containers/%s/checkpoints/", srcAddr, containerId, dstAddr, dstContainerId)

	err = cmd.Run()
	if err != nil {
		return err
	}

	// restore execution on the new server
	err = cliDst.ContainerStart(context.Background(), dstContainerId, types.ContainerStartOptions{CheckpointID: "predumpCheckpointB"})
	if err != nil {
		return err
	}

	return nil
}

// one-pass migrate,
// freezes the container and migrate
func (daemon *Daemon) migrate(containerId string, cliSrc *client.Client, cliDst *client.Client, srcAddr string, dstAddr string) error {
	done := make(chan int, 10)

	// find the docker image and configs
	inspectJson, err := cliSrc.ContainerInspect(context.Background(),
		containerId)
	if err != nil {
		return err
	}

	// create a container in the destination container
	createResp, dstErr := cliDst.ContainerCreate(context.Background(),
		inspectJson.Config,
		inspectJson.HostConfig,
		&network.NetworkingConfig{EndpointsConfig: inspectJson.NetworkSettings.Networks},
		&specs.Platform{
			OS:           "linux",
			Architecture: "amd64",
		},
		containerId)
	if dstErr != nil {
		return dstErr
	}

	dstContainerId := createResp.ID

	// sync
	<-done

	// second time checkpoint the source container
	srcErr := cliSrc.CheckpointCreate(context.Background(),
		containerId,
		types.CheckpointCreateOptions{
			Predump:      false,
			Exit:         true,
			CheckpointID: "checkpointA"})
	if srcErr != nil {
		return srcErr
	}

	// scp the checkpoints to the new server
	cmd := exec.Command("scp %s:/var/lib/docker/containers/%s/checkpoints/checkpointA "+
		"%s:/var/lib/docker/containers/%s/checkpoints/", srcAddr, containerId, dstAddr, dstContainerId)

	err = cmd.Run()
	if err != nil {
		return err
	}

	// restore execution on the new server
	err = cliDst.ContainerStart(context.Background(), dstContainerId, types.ContainerStartOptions{CheckpointID: "checkpointA"})
	if err != nil {
		return err
	}

	return nil
}
