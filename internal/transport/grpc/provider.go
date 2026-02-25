package grpc

import (
	"go.uber.org/fx"
)

var Module = fx.Module("grpc",
	fx.Provide(NewServer),
	fx.Invoke(RegisterLifecycle),
)
