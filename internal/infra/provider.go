package infra

import (
	"go.uber.org/fx"

	"cms-api/internal/infra/cache"
	"cms-api/internal/infra/database"
	"cms-api/internal/infra/httpclient"
	"cms-api/internal/infra/search"
	"cms-api/internal/infra/telemetry"
)

var Module = fx.Module("infra",
	database.Module,
	httpclient.Module,
	search.Module,
	telemetry.Module,
	cache.Module,
)
