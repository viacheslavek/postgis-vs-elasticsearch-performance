package storage

import "context"

type Storage interface {
	AddPoint(ctx context.Context, p Point) error
	AddBatch(ctx context.Context, points []Point) error
	DeletePoint(ctx context.Context, p Point) error
	IsPoint(ctx context.Context, p Point) (bool, error)

	GetInRadius(p Point, radius float64) ([]Point, error)
	GetInPolygon(polygon []float64) ([]Point, error)
}

type Point struct {
	Latitude  float64
	Longitude float64
}
