package discovery

import (
	"go.uber.org/fx"

	discoveryhttp "cms-api/internal/modules/discovery/http"
	"cms-api/internal/modules/discovery/repo"
	"cms-api/internal/modules/discovery/service"
)

var Module = fx.Module("discovery",
	fx.Provide(repo.New),
	fx.Provide(service.New),
	fx.Provide(discoveryhttp.NewHandler),
	fx.Invoke(discoveryhttp.RegisterRoutes),
)
