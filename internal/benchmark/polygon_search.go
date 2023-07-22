package benchmark

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"time"
)

func benchGetInRadiusPolygon(ctx context.Context, p internal.Polygon, radius int) (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func benchGetInPolygonPolygon(ctx context.Context, polygon internal.Polygon) (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func benchGetIntersectionPolygon(ctx context.Context, polygon internal.Polygon) (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func benchGetIntersectionPoint(ctx context.Context, point internal.Point) (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}
