package auth

import (
	"crypto/rsa"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"cms-api/internal/config"
	authhttp "cms-api/internal/modules/auth/http"
	"cms-api/internal/modules/auth/repo"
	"cms-api/internal/modules/auth/service"
	"cms-api/internal/pkg/crypto"
)

func newService(r repo.Repository, cfg *config.Config, log *zap.Logger) (service.Service, error) {
	var key *rsa.PrivateKey
	if cfg.JWT.PrivateKeyPath != "" {
		var err error
		key, err = crypto.LoadPrivateKey(cfg.JWT.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
	}
	return service.New(r, key, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry, log), nil
}

var Module = fx.Module("auth",
	fx.Provide(repo.New),
	fx.Provide(newService),
	fx.Provide(authhttp.NewHandler),
	fx.Invoke(authhttp.RegisterRoutes),
)
