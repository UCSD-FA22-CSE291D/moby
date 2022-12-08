package migrate

import "github.com/docker/docker/api/types"

type Backend interface {
	Migrate(containerId string, options types.MigrateOptions) error
}
