package benchmark

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"time"
)

func benchInitPolygon(ctx context.Context, s storage.Storage) (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func benchDropPolygon(ctx context.Context, s storage.Storage) (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func benchAddPolygon(ctx context.Context, s storage.Storage) (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func benchAddPolygonBatch(ctx context.Context, s storage.Storage) (time.Duration, error) {
	// новая стратегия - генерировать точки внутри
	//TODO implement me
	panic("implement me")
}
