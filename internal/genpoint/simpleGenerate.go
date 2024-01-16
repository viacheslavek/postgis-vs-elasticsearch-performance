package genpoint

import (
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"math"
	"math/rand"
)

type SimplePointGenerator struct{}

// GeneratePoints Задает N точек в окружности в пределах Московской области
func (smg *SimplePointGenerator) GeneratePoints(N int) []internal.Point {

	centerMosRegion := internal.Point{
		Latitude:  55.751426,
		Longitude: 37.618879,
	}

	edgeMosRegion := internal.Point{
		Latitude:  56.342905,
		Longitude: 37.517608,
	}

	points := make([]internal.Point, N)

	radiusX := int(math.Abs(centerMosRegion.Latitude-edgeMosRegion.Latitude) * 10e6)
	radiusY := int(math.Abs(centerMosRegion.Longitude-edgeMosRegion.Longitude) * 10e6)

	for i := 0; i < N; i++ {
		points[i] = generatePointInRadius(
			int(centerMosRegion.Latitude*10e6),
			int(centerMosRegion.Longitude*10e6),
			radiusX, radiusY)
	}

	return points
}

func generatePointInRadius(centralX, centralY, radiusX, radiusY int) internal.Point {

	newRadiusX := rand.Intn(radiusX)
	newRadiusY := rand.Intn(radiusY)

	xRand := centralX
	yRand := centralY

	signX := rand.Intn(2)
	if signX == 0 {
		xRand += newRadiusX
	} else {
		xRand -= newRadiusX
	}
	signY := rand.Intn(2)
	if signY == 0 {
		yRand += newRadiusY
	} else {
		yRand += newRadiusY
	}

	return internal.Point{
		Latitude:  float64(xRand/10) / 10e5,
		Longitude: float64(yRand/10) / 10e5,
	}
}
