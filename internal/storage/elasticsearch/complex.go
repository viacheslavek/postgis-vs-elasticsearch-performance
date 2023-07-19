package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"io"
	"log"
	"strconv"
	"strings"
)

type Response struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Point string `json:"point"`
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

	return s.doRequestWithQuery(ctx, qJSON)
}

func (s *Storage) GetInPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Point, error) {

	qJson := fmt.Sprintf(`
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

	return s.doRequestWithQuery(ctx, qJson)
}

func (s *Storage) doRequestWithQuery(ctx context.Context, qJSON string) ([]internal.Point, error) {
	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex("moscow_region"),
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

	return getInternalPointsFromStructString(response), nil
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
