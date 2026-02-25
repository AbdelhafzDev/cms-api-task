package youtube

import (
	"context"
	"time"

	"go.uber.org/zap"

	"cms-api/internal/modules/importer/service"
)

type youtubeStub struct {
	log *zap.Logger
}

func NewYouTubeStubImporter(log *zap.Logger) *youtubeStub {
	return &youtubeStub{log: log.Named("youtube_importer")}
}

func (y *youtubeStub) SourceType() string {
	return "youtube"
}

func (y *youtubeStub) Fetch(ctx context.Context, baseURL string, since *time.Time) ([]service.ImportItem, error) {
	_ = ctx
	y.log.Info("YouTube importer stub invoked", zap.String("base_url", baseURL))
	if since != nil {
		y.log.Info("YouTube importer stub since", zap.Time("since", *since))
	}
	return []service.ImportItem{}, nil
}
