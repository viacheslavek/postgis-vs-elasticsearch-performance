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

type PolygonStorage interface {
	InitPolygon(ctx context.Context) error
	DropPolygon(ctx context.Context) error

	AddPolygon(ctx context.Context, polygon internal.Polygon) error
	AddPolygonBatch(ctx context.Context, polygons []internal.Polygon) error

	GetInRadiusPolygon(ctx context.Context, p internal.Polygon, radius int) ([]internal.Polygon, error)
	GetInPolygonPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Polygon, error)

	GetIntersectionPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Polygon, error)
	GetIntersectionPoint(ctx context.Context, point internal.Point) ([]internal.Polygon, error)
}
