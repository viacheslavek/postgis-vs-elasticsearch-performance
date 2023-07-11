package storage

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
)

type Storage interface {
	AddPoint(ctx context.Context, p internal.Point) error
	AddBatch(ctx context.Context, points []internal.Point) error
	DeletePoint(ctx context.Context, p internal.Point) error
	IsPoint(ctx context.Context, p internal.Point) (bool, error)

	GetInRadius(p internal.Point, radius float64) ([]internal.Point, error)
	GetInPolygon(polygon []float64) ([]internal.Point, error)
}
