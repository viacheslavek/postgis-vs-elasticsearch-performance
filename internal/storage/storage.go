package storage

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
)

type Storage interface {
	Init(ctx context.Context) error
	Drop(ctx context.Context) error

	AddPoint(ctx context.Context, p internal.Point) error
	AddPointBatch(ctx context.Context, points []internal.Point) error

	GetInRadius(ctx context.Context, p internal.Point, radius int) ([]internal.Point, error)
	GetInPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Point, error)
}
