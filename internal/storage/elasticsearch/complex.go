package elasticsearch

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
)

func (s *Storage) GetInRadius(ctx context.Context, p internal.Point, radius int) ([]internal.Point, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetInPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Point, error) {
	//TODO implement me
	panic("implement me")
}
