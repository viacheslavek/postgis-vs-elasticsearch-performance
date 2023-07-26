package benchmark

const (
	PointInit     = "time to init point db: "
	PointDrop     = "time to drop point db: "
	PointAddBatch = "time to add batch points in db: "

	PointSearchInRadius = "time to search points in radius: "
	PointSearchInShapes = "time to search points in shapes: "

	PolygonInit     = "time to init polygon db: "
	PolygonDrop     = "time to drop polygon db: "
	PolygonAddBatch = "time to add batch polygon in db: "

	PolygonSearchInRadius       = "time to search polygons in radius: "
	PolygonSearchInPolygon      = "time to search polygons in polygon: "
	PolygonGetIntersection      = "time to search polygons intersection with polygon: "
	PolygonGetIntersectionPoint = "time to search polygons intersection with point: "
)
