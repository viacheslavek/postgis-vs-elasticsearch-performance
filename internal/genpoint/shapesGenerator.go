package genpoint

import (
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"math/rand"
)

func GenerateShapes(count int) internal.Shapes {
	countPolygon := rand.Intn(count)
	countCircle := count - countPolygon

	polGen := PolygonGenerator{}

	shapes := internal.Shapes{
		Polygons: polGen.GeneratePolygons(countPolygon),
		Circles:  generateCircles(countCircle),
	}

	return shapes
}

func generateCircles(count int) []internal.Circle {
	circles := make([]internal.Circle, count)
	spg := SimplePointGenerator{}
	for i := 0; i < count; i++ {
		circles[i] = generateCircle(spg)
	}

	return circles
}

func generateCircle(spg SimplePointGenerator) internal.Circle {
	return internal.Circle{
		Centre: spg.GeneratePoints(1)[0],
		Radius: rand.Intn(5e4),
	}
}
