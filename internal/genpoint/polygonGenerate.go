package genpoint

import (
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"math"
	"sort"
)

type PolygonGenerator struct{}

type ByAngel struct {
	Points    []internal.Point
	Center    internal.Point
	Distances []float64
}

func (a ByAngel) Len() int {
	return len(a.Points)
}
func (a ByAngel) Less(i, j int) bool {
	return a.Distances[i] < a.Distances[j]
}
func (a ByAngel) Swap(i, j int) {
	a.Points[i], a.Points[j] = a.Points[j], a.Points[i]
	a.Distances[i], a.Distances[j] = a.Distances[j], a.Distances[i]
}

func (pg *PolygonGenerator) GeneratePolygon(N int) []internal.Point {
	pointsGen := SimplePointGenerator{}
	points := pointsGen.GeneratePoints(N)

	// -> перенести этот закомментированный код в тесты
	//points := []internal.Point{
	//	{1, 1},
	//	{1, -1},
	//	{-1, 1},
	//	{-1, -1},
	//}

	// алгоритм:
	// Находим точку в центре многоугольника
	// сортируем все точки по арктангенсу относительно центральной точки

	centre := centrePoint(points)

	sort.Sort(ByAngel{
		Points:    points,
		Center:    centre,
		Distances: getDistances(centre, points),
	})

	return points
}

func centrePoint(points []internal.Point) internal.Point {
	sumX, sumY := 0.0, 0.0

	for _, point := range points {
		sumX += point.Latitude
		sumY += point.Longitude
	}

	pX := math.Round(sumX/float64(len(points))*10e6) / 10e6
	pY := math.Round(sumY/float64(len(points))*10e6) / 10e6

	return internal.Point{
		Latitude:  pX,
		Longitude: pY,
	}
}

func getDistances(centre internal.Point, points []internal.Point) []float64 {
	distances := make([]float64, len(points))

	for i, p := range points {
		distances[i] = math.Atan2(p.Longitude-centre.Longitude, p.Latitude-centre.Latitude)
	}

	return distances
}
