package daemon // import "github.com/docker/docker/daemon"

import (
	"context"
	"fmt"
	// "os"
	// "path/filepath"

	"github.com/docker/docker/api/types"
	// "github.com/docker/docker/daemon/names"
	// "github.com/docker/docker/daemon/checkpoint"
)

func (daemon *Daemon) Migrate(name string, config types.MigrateOptions) error {
	//TODO: may add pre-dumps before the final dump

	checkpointConfig := types.CheckpointCreateOptions {
		CheckpointID:		"MigrateFinalDump",
		CheckpointDir:		"",
		Predump:			false,
		ParentCheckpointID: "",
		Exit:				true,
	}

	err := daemon.CheckpointCreate(name, checkpointConfig)
	if err != nil {
		return err
	}

	targetAddr := config.TargetAddr
	if targetAddr != "localhost" { //TODO: support migration with remote dest
		return fmt.Errorf("Only support on migration to localhost from container %s", name)
	}

	//TODO: run the new instance with dumped files
	// (used with: docker start --checkpoint CHECKPOINT_ID [OTHER OPTIONS] CONTAINER)
	err = daemon.ContainerStart(context.Background(), name, nil, checkpointConfig.CheckpointID, checkpointConfig.CheckpointDir)
	if err != nil {
		return err
	}

	return nil
}
