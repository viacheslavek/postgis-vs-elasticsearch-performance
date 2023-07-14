package elasticsearch

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"log"
	"strings"
)

// получить все индексы curl -XGET 'http://localhost:9200/_cat/indices?v'

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

	if exist, err := s.isExist(ctx); err != nil || exist {
		if err != nil {
			return fmt.Errorf("can't check inExist %w\n", err)
		}
		return nil
	}

	indexRequest := esapi.IndicesCreateRequest{
		Index: "moscow_region",
		Body: strings.NewReader(`
			{
				"mappings": {
					"properties": {
						"id": {
							"type": "integer"
						},
						"point": {
							"type": "geo_point"
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

func (s *Storage) Drop(ctx context.Context) error {

	if exist, err := s.isExist(ctx); err != nil || !exist {
		if err != nil {
			return fmt.Errorf("can't check inExist %w\n", err)
		}
		return nil
	}

	deleteIndexRequest := esapi.IndicesDeleteRequest{
		Index: []string{"moscow_region"},
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

func (s *Storage) isExist(ctx context.Context) (bool, error) {
	indexExistsRequest := esapi.IndicesExistsRequest{
		Index: []string{"moscow_region"},
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
