package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/genpoint"
	"io"
	"log"
	"strconv"
	"strings"
)

type Response struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Point   string `json:"point"`
				Polygon struct {
					Coordinates [][][]float64 `json:"coordinates"`
				} `json:"polygon"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (s *Storage) GetInRadius(ctx context.Context, p internal.Point, radius int) ([]internal.Point, error) {

	qJSON := fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": {
					"geo_distance": {
						"distance": "%dm",
						"point": "%s"
					}
				}
			}
		},
		"size": "1000"
	}
	`, radius, convertPointToES(p))

	resp, err := s.doSearchRequestWithQuery(ctx, qJSON, "moscow_region")

	if err != nil {
		return nil, fmt.Errorf("can't doSearchPointRequestWithQuery InRadius %w\n", err)
	}

	return getInternalPointsFromStructString(*resp), nil
}

func (s *Storage) GetInPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Point, error) {

	qJSON := fmt.Sprintf(`
	{
		"query": {
			"geo_polygon": {
				"point": {
					"points": [%s]
				}
			}
		},
		"size": "1000"
	}`, generateESPolygon(polygon))

	resp, err := s.doSearchRequestWithQuery(ctx, qJSON, "moscow_region")

	if err != nil {
		return nil, fmt.Errorf("can't doSearchPointRequestWithQuery InPolygon %w\n", err)
	}

	return getInternalPointsFromStructString(*resp), nil
}

func (s *Storage) GetInRadiusPolygon(ctx context.Context, p internal.Polygon, radius int) ([]internal.Polygon, error) {

	centralPoint, newRadius := genpoint.GetCentralPolygonPointWithRadius(p.Vertical, radius)

	qJSON := fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": {
					"geo_shape": {
						"polygon": {
							"shape": {
								"type":"circle",
								"coordinates": [%f, %f],
								"radius": "%dm"
							},
							"relation": "within"
						}
					}
				}
			}
		},
		"size": "1000"
	}`, centralPoint.Latitude, centralPoint.Longitude, newRadius)

	resp, err := s.doSearchRequestWithQuery(ctx, qJSON, "moscow_region_polygon")

	if err != nil {
		return nil, fmt.Errorf("can't doSearchPointRequestWithQuery InRadiusPolygon %w\n", err)
	}

	return getInternalPolygonsFromStructString(*resp), nil

}

func (s *Storage) GetInPolygonPolygon(ctx context.Context, polygon internal.Polygon) ([]internal.Polygon, error) {

	qJSON := fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": {
					"geo_shape": {
						"polygon": {
							"relation": "within",
							"shape": {
								"type": "polygon",
								"coordinates": [
									[%s]
								]
							}
						}
					}
				}
			}
		},
		"size": "10000"
	}`, generateESPolygon(polygon.Vertical))

	resp, err := s.doSearchRequestWithQuery(ctx, qJSON, "moscow_region_polygon")

	if err != nil {
		return nil, fmt.Errorf("can't doSearchPointRequestWithQuery InPolygonPolygon %w\n", err)
	}

	return getInternalPolygonsFromStructString(*resp), nil
}

func (s *Storage) GetIntersectionPolygon(ctx context.Context, polygon internal.Polygon) ([]internal.Polygon, error) {
	qJSON := fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": {
					"geo_shape": {
						"polygon": {
							"relation": "intersects",
							"shape": {
								"type": "polygon",
								"coordinates": [
									[%s]
								]
							}
						}
					}
				}
			}
		},
		"size": "10000"
	}`, generateESPolygon(polygon.Vertical))

	resp, err := s.doSearchRequestWithQuery(ctx, qJSON, "moscow_region_polygon")

	if err != nil {
		return nil, fmt.Errorf("can't doSearchPointRequestWithQuery InPolygonPolygon %w\n", err)
	}

	return getInternalPolygonsFromStructString(*resp), nil
}

func (s *Storage) GetIntersectionPoint(ctx context.Context, point internal.Point) ([]internal.Polygon, error) {
	qJSON := fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": {
					"geo_shape": {
						"polygon": {
							"relation": "contains",
							"shape": {
								"type": "point",
								"coordinates": [%f, %f]
							}
						}
					}
				}
			}
		},
		"size": "10000"
	}`, point.Latitude, point.Longitude)

	resp, err := s.doSearchRequestWithQuery(ctx, qJSON, "moscow_region_polygon")

	if err != nil {
		return nil, fmt.Errorf("can't doSearchPointRequestWithQuery InPolygonPolygon %w\n", err)
	}

	return getInternalPolygonsFromStructString(*resp), nil
}

func (s *Storage) GetInShapes(ctx context.Context, shape internal.Shapes) ([]internal.Point, error) {

	qJSON := fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": {
					"geo_shape": {
						"point": {
							"relation": "intersects",
							"shape": {
								"type": "geometrycollection",
								"geometries": [
									%s
								]
							}
						}
					}
				}
			}
		},
		"size": "10000"
	}`, convertShapeToEsStr(shape))

	resp, err := s.doSearchRequestWithQuery(ctx, qJSON, "moscow_region")

	if err != nil {
		return nil, fmt.Errorf("can't doSearchPointRequestWithQuery InPolygon %w\n", err)
	}

	return getInternalPointsFromStructString(*resp), nil
}

func (s *Storage) doSearchRequestWithQuery(ctx context.Context, qJSON, index string) (*Response, error) {
	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(index),
		s.client.Search.WithBody(strings.NewReader(qJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("can't do request to search %w\n", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("can't close body exist request\n")
		}
	}(res.Body)

	if res.IsError() {
		return nil, fmt.Errorf("can't do request exist %s\n", res.String())
	}

	var buf bytes.Buffer

	if _, err = io.Copy(&buf, res.Body); err != nil {
		return nil, fmt.Errorf("can't copy body request %w", err)
	}

	var response Response

	if err = json.Unmarshal(buf.Bytes(), &response); err != nil {
		return nil, fmt.Errorf("can't unmarshal Json %w\n", err)
	}

	return &response, nil
}

func convertShapeToEsStr(shapes internal.Shapes) string {

	lenPolygon := len(shapes.Polygons)
	lenCircle := len(shapes.Circles)

	queries := make([]string, lenPolygon+lenCircle)

	for i := 0; i < lenPolygon; i++ {
		queries[i] = convertPolygonToES(shapes.Polygons[i])
	}

	for i := 0; i < lenCircle; i++ {
		queries[i+lenPolygon] = convertCircleToES(shapes.Circles[i])
	}

	return strings.Join(queries, ", ")
}

func convertPolygonToES(polygon internal.Polygon) string {

	return fmt.Sprintf(`
		{
			"type": "polygon",
			"coordinates": [[%s]]
		}`, generateESPolygon(polygon.Vertical))
}

func convertCircleToES(circle internal.Circle) string {
	return fmt.Sprintf(`
		{
			"type": "circle",
			"coordinates": [%f, %f],
	 		"radius": "%dm"
		}`, circle.Centre.Longitude, circle.Centre.Latitude, circle.Radius)
}

func getInternalPolygonsFromStructString(response Response) []internal.Polygon {
	polygons := make([]internal.Polygon, 0)
	for _, hit := range response.Hits.Hits {
		polygon := internal.Polygon{
			Vertical: getPointsFromFloat(hit.Source.Polygon.Coordinates),
		}
		polygons = append(polygons, polygon)
	}
	return polygons
}

func getPointsFromFloat(coordinate [][][]float64) []internal.Point {

	points := make([]internal.Point, len(coordinate[0]))
	for i, p := range coordinate[0] {
		points[i].Latitude = p[0]
		points[i].Longitude = p[1]
	}
	return points
}

func getInternalPointsFromStructString(response Response) []internal.Point {
	points := make([]internal.Point, 0)
	for _, hit := range response.Hits.Hits {
		points = append(points, convertStrToInternal(hit.Source.Point))
	}
	return points
}

func convertStrToInternal(point string) internal.Point {
	latLon := strings.Split(point, ", ")

	lat, err := strconv.ParseFloat(latLon[0], 64)
	if err != nil {
		log.Printf("can't convert es point to internal float")
	}
	lon, err := strconv.ParseFloat(latLon[1], 64)
	if err != nil {
		log.Printf("can't convert es point to internal float")
	}

	return internal.Point{
		Latitude:  lat,
		Longitude: lon,
	}
}

func generateESPolygon(points []internal.Point) string {
	esPoints := &strings.Builder{}

	N := len(points)
	for i, point := range points {
		_, err := fmt.Fprintf(esPoints, "[%f, %f]", point.Latitude, point.Longitude)
		if err != nil {
			log.Println("can't convert internalP to ESP")
		}
		if i < N-1 {
			esPoints.WriteByte(',')
			esPoints.WriteByte(' ')
		}
	}

	return esPoints.String()
}
