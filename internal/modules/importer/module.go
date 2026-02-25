package importer

import (
	"go.uber.org/fx"

	"cms-api/internal/modules/importer/adapter/youtube"
	importhttp "cms-api/internal/modules/importer/http"
	"cms-api/internal/modules/importer/repo"
	"cms-api/internal/modules/importer/service"
)

var Module = fx.Module("importer",
	fx.Provide(repo.New),
	fx.Provide(service.NewRegistry),
	fx.Provide(service.New),
	fx.Provide(
		fx.Annotate(
			youtube.NewYouTubeStubImporter,
			fx.ResultTags(`group:"importers"`),
		),
	),
	fx.Provide(importhttp.NewHandler),
	fx.Invoke(importhttp.RegisterRoutes),
)
