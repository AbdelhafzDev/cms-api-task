package transport

import (
	"go.uber.org/fx"

	"cms-api/internal/transport/grpc"
	"cms-api/internal/transport/http"
)

var Module = fx.Module("transport",
	http.Module,
	grpc.Module,
)
