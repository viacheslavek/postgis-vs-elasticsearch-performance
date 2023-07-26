package benchmark

const (
	PointInit = iota
	PointDrop
	PointAddBatch

	PointSearchInRadius
	PointSearchInShapes

	PolygonInit
	PolygonDrop
	PolygonAddBatch

	PolygonSearchInRadius
	PolygonSearchInPolygon
	PolygonGetIntersection
	PolygonGetIntersectionPoint
)
