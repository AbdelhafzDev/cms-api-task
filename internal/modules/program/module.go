package program

import (
	"go.uber.org/fx"

	programhttp "cms-api/internal/modules/program/http"
	"cms-api/internal/modules/program/repo"
	"cms-api/internal/modules/program/service"
)

var Module = fx.Module("program",
	fx.Provide(repo.New),
	fx.Provide(service.New),
	fx.Provide(programhttp.NewHandler),
	fx.Invoke(programhttp.RegisterRoutes),
)
