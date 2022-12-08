package migrate

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
)

func (s *migrateRouter) postContainerMigrate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	var options types.MigrateOptions
	if err := httputils.ReadJSON(r, &options); err != nil {
		return err
	}

	err := s.backend.Migrate(vars["id"], options)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
