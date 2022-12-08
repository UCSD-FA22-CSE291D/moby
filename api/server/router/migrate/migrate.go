package migrate

import (
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/router"
)

// migrateRouter
type migrateRouter struct {
	backend Backend
	decoder httputils.ContainerDecoder
	routes  []router.Route
}

// NewRouter initializes a new checkpoint router
func NewRouter(b Backend, decoder httputils.ContainerDecoder) router.Router {
	r := &migrateRouter{
		backend: b,
		decoder: decoder,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the checkpoint controller
func (s *migrateRouter) Routes() []router.Route {
	return s.routes
}

func (s *migrateRouter) initRoutes() {
	s.routes = []router.Route{
		router.NewPostRoute("/containers/{id:.*}/migrate", s.postContainerMigrate, router.Experimental),
	}
}
