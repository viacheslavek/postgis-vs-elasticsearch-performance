package elasticsearch

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"
)

type Storage struct {
	client *elasticsearch.Client
}

func New() (*Storage, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create new client to es %w\n", err)
	}

	return &Storage{client: es}, nil
}

func (s *Storage) Ping() error {

	req := esapi.InfoRequest{}

	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		return fmt.Errorf("can't do req to es %w", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("Error response: %s\n", res.String())
	}

	fmt.Println(res.String())
	return nil
}

func (s *Storage) Init(ctx context.Context) error {

	mapping := `
		{
			"mappings": {
				"properties": {
					"id": {
						"type": "keyword"
					},
					"point": {
						"type": "geo_point"
					}
				}
			}
		}
	`

	return s.initIndex(ctx, "moscow_region", mapping)
}

func (s *Storage) Drop(ctx context.Context) error {
	return s.drop(ctx, "moscow_region")
}

func (s *Storage) AddPoint(ctx context.Context, p internal.Point) error {

	doc := getDocument(convertPointToES(p))

	indexRequest := esapi.IndexRequest{
		Index:      "moscow_region",
		DocumentID: doc.Id,
		Body: strings.NewReader(
			fmt.Sprintf(`{"point": "%s"}`, doc.GeoPoint)),
	}

	return s.addSingle(ctx, indexRequest)
}

func (s *Storage) AddPointBatch(ctx context.Context, points []internal.Point) error {

	pointsES := convertPointsToES(points)

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: s.client,
		Index:  "moscow_region",
	})
	if err != nil {
		return fmt.Errorf("can't init bulk connect %w\n", err)
	}

	for _, point := range pointsES {
		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: point.Id,
			Body: strings.NewReader(
				fmt.Sprintf(`{"point": "%s"}`, point.GeoPoint)),
		})
		if err != nil {
			return fmt.Errorf("can't add point to bulk indexer %w\n", err)
		}
	}

	err = bi.Close(ctx)
	if err != nil {
		return fmt.Errorf("can;t close bulk %w\n", err)
	}

	return nil
}

type esPoint struct {
	Id       string `json:"id"`
	GeoPoint string `json:"geo_point"`
}

func (s *Storage) initIndex(ctx context.Context, index, mapping string) error {
	if exist, err := s.isExist(ctx, index); err != nil || exist {
		if err != nil {
			return fmt.Errorf("can't check inExist %w\n", err)
		}
		return nil
	}

	indexRequest := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(mapping),
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

	return nil

}

func (s *Storage) drop(ctx context.Context, index string) error {
	if exist, err := s.isExist(ctx, index); err != nil || !exist {
		if err != nil {
			return fmt.Errorf("can't check inExist %w\n", err)
		}
		return nil
	}

	deleteIndexRequest := esapi.IndicesDeleteRequest{
		Index: []string{index},
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

func (s *Storage) addSingle(ctx context.Context, req esapi.IndexRequest) error {
	indexResponse, err := req.Do(ctx, s.client)
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

func (s *Storage) isExist(ctx context.Context, index string) (bool, error) {
	indexExistsRequest := esapi.IndicesExistsRequest{
		Index: []string{index},
	}

	indexExistsResponse, err := indexExistsRequest.Do(ctx, s.client)
	if err != nil {
		return false, fmt.Errorf("Error checking if index exists: %w\n", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("can`t close body in ExistRequest\n")
		}
	}(indexExistsResponse.Body)

	if !indexExistsResponse.IsError() {
		fmt.Println("Index already exists")
		return true, nil
	}

	return false, nil
}

func convertPointToES(p internal.Point) string {
	return fmt.Sprintf("%f, %f", p.Longitude, p.Latitude)
}

func convertPointsToES(points []internal.Point) []esPoint {
	esPoints := make([]esPoint, len(points))
	for i, p := range points {
		esPoints[i] = getDocument(convertPointToES(p))
	}

	return esPoints
}

func getDocument(pointES string) esPoint {
	return esPoint{
		Id:       uuid.New().String(),
		GeoPoint: pointES,
	}
}
