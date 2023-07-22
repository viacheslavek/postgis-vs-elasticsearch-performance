package genpoint

import (
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"math"
	"math/rand"
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

	points = append(points, points[0])

	return points
}

func (pg *PolygonGenerator) GeneratePolygons(N int) []internal.Polygon {

	polygons := make([]internal.Polygon, N)

	for i := 0; i < N; i++ {
		angel := rand.Intn(30)
		// get rid of zero-two
		angel += 3
		polygons[i] = internal.Polygon{Vertical: pg.GeneratePolygon(angel)}
	}

	return polygons

}

func GetCentralPolygonPointWithRadius(polygon []internal.Point, radius int) (internal.Point, int) {

	centre := centrePoint(polygon)

	centre.Longitude = math.Round(centre.Longitude*10e5) / 10e5
	centre.Latitude = math.Round(centre.Latitude*10e5) / 10e5

	newRadius := radius + getDelta(centre, polygon)

	return centre, newRadius
}

func centrePoint(points []internal.Point) internal.Point {
	sumX, sumY := 0.0, 0.0

	for _, point := range points {
		sumX += point.Latitude
		sumY += point.Longitude
	}

	pX := math.Round(sumX/float64(len(points))*10e5) / 10e5
	pY := math.Round(sumY/float64(len(points))*10e5) / 10e5

	return internal.Point{
		Latitude:  pX,
		Longitude: pY,
	}
}

func getDelta(centre internal.Point, polygon []internal.Point) int {
	centreOrb := orb.Point{centre.Longitude, centre.Latitude}

	polygonOrb := make([]orb.Point, len(polygon))

	for i, p := range polygon {
		polygonOrb[i] = orb.Point{p.Longitude, p.Latitude}
	}

	var maxDistance float64
	for _, point := range polygonOrb {
		distance := geo.DistanceHaversine(centreOrb, point)
		maxDistance = math.Max(maxDistance, distance)
	}

	return int(math.Round(maxDistance))

}

func getDistances(centre internal.Point, points []internal.Point) []float64 {
	distances := make([]float64, len(points))

	for i, p := range points {
		distances[i] = math.Atan2(p.Longitude-centre.Longitude, p.Latitude-centre.Latitude)
	}

	return distances
}
