package app

import (
	"go.uber.org/fx"

	"cms-api/internal/modules/auth"
	"cms-api/internal/modules/discovery"
	"cms-api/internal/modules/importer"
	"cms-api/internal/modules/program"
	"cms-api/internal/modules/worker"
)

var FeatureModules = fx.Options(
	auth.Module,
	worker.Module,
	program.Module,
	discovery.Module,
	importer.Module,
)
