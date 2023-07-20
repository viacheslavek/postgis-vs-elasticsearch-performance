package elasticsearch

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"
)

// TODO: отрефакторить код так, чтобы я для всех запросов в es передавал только queryJSON в вспомогательную функцию,
// тем самым сокращу код

func (s *Storage) InitPolygon(ctx context.Context) error {

	if exist, err := s.isExist(ctx, "moscow_region_polygon"); err != nil || exist {
		if err != nil {
			return fmt.Errorf("can't check inExist %w\n", err)
		}
		return nil
	}

	indexRequest := esapi.IndicesCreateRequest{
		Index: "moscow_region_polygon",
		Body: strings.NewReader(`
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
		`),
	}

	indexResponse, err := indexRequest.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("can't do request mapping %w\n", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("can`t close body\n")
		}
	}(indexResponse.Body)

	if indexResponse.IsError() {
		return fmt.Errorf("index response have err %v\n", indexResponse.String())
	}

	log.Println("index created success")

	return nil
}

func (s *Storage) DropPolygon(ctx context.Context) error {

	if exist, err := s.isExist(ctx, "moscow_region_polygon"); err != nil || !exist {
		if err != nil {
			return fmt.Errorf("can't check inExist %w\n", err)
		}
		return nil
	}

	deleteIndexRequest := esapi.IndicesDeleteRequest{
		Index: []string{"moscow_region_polygon"},
	}

	deleteIndexResponse, err := deleteIndexRequest.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("can't do delete request %w\n", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("can't close body in delete\n")
		}
	}(deleteIndexResponse.Body)

	if deleteIndexResponse.IsError() {
		return fmt.Errorf("Error deleting the index: %s\n", deleteIndexResponse.String())
	}

	return nil
}

func (s *Storage) AddPolygon(ctx context.Context, p internal.Polygon) error {

	docId := uuid.New().String()

	indexRequest := esapi.IndexRequest{
		Index:      "moscow_region_polygon",
		DocumentID: docId,
		Body: strings.NewReader(
			fmt.Sprintf(`
				{"geo_polygon": {
					"points": [%s]
				}}`, generateESPolygon(p.Vertical))),
	}

	fmt.Println(indexRequest.Body)

	indexResponse, err := indexRequest.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("can't do request to add point %w\n", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("can`t close body add reqest\n")
		}
	}(indexResponse.Body)

	if indexResponse.IsError() {
		return fmt.Errorf("Error add the index: %s\n", indexResponse.String())
	}

	return nil
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
				fmt.Sprintf(`{"geo_polygon": {"points": [%s]}}`, poly)),
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

func (s *Storage) GetInRadiusPolygon(ctx context.Context, p internal.Polygon, radius int) ([]internal.Polygon, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetInPolygonPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Polygon, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetIntersectionPolygon(ctx context.Context, polygon []internal.Point) ([]internal.Polygon, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetIntersectionPoint(ctx context.Context, point internal.Point) ([]internal.Polygon, error) {
	//TODO implement me
	panic("implement me")
}

func convertPolygonsToES(polygons []internal.Polygon) []string {
	polygonsES := make([]string, len(polygons))

	for i, p := range polygons {
		polygonsES[i] = generateESPolygon(p.Vertical)
	}
	return polygonsES
}
