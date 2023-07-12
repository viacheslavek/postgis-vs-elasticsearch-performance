package genpoint

import "github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"

type GetPoint interface {
	GeneratePoints(N int) []internal.Point
}

type GetPolygon interface {
	GeneratePolygon(N int) []internal.Point
}
