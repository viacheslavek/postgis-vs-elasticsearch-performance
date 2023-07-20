package elasticsearch

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"
	"strings"
)

func (s *Storage) InitPolygon(ctx context.Context) error {

	mapping := `
		{
			"mappings": {
				"properties": {
					"id": {
						"type": "keyword"
					},
					"polygon": {
						"type": "geo_shape"
					}
				}
			}
		}	
	`

	return s.initIndex(ctx, "moscow_region_polygon", mapping)
}

func (s *Storage) DropPolygon(ctx context.Context) error {
	return s.drop(ctx, "moscow_region_polygon")
}

func (s *Storage) AddPolygon(ctx context.Context, p internal.Polygon) error {

	indexRequest := esapi.IndexRequest{
		Index:      "moscow_region_polygon",
		DocumentID: uuid.New().String(),
		Body: strings.NewReader(
			fmt.Sprintf(`
				{"polygon": {
					"type": "polygon",
					"coordinates": [[%s]]
				}}`, generateESPolygon(p.Vertical))),
	}

	return s.addSingle(ctx, indexRequest)
}

func (s *Storage) AddPolygonBatch(ctx context.Context, polygon []internal.Polygon) error {

	polygonES := convertPolygonsToES(polygon)

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: s.client,
		Index:  "moscow_region_polygon",
	})
	if err != nil {
		return fmt.Errorf("can't init bulk connect %w\n", err)
	}

	for _, poly := range polygonES {
		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: uuid.New().String(),
			Body: strings.NewReader(
				fmt.Sprintf(`{"polygon": {"type": "polygon", "coordinates": [[%s]] }}`, poly)),
		})

		if err != nil {
			return fmt.Errorf("can't add polygon to bulk indexer %w\n", err)
		}
	}

	err = bi.Close(ctx)
	if err != nil {
		return fmt.Errorf("can`t close bulk %w\n", err)
	}

	return nil
}

func convertPolygonsToES(polygons []internal.Polygon) []string {
	polygonsES := make([]string, len(polygons))

	for i, p := range polygons {
		polygonsES[i] = generateESPolygon(p.Vertical)
	}
	return polygonsES
}
