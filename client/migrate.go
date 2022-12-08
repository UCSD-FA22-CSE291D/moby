package client

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (cli *Client) Migrate(ctx context.Context, containerId string, options types.MigrateOptions) error {
	resp, err := cli.post(ctx, "/containers/"+containerId+"/migrate", nil, options, nil)
	ensureReaderClosed(resp)
	return err
}
