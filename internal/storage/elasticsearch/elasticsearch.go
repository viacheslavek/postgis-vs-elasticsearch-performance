package elasticsearch

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
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

	// Вывод информации о кластере
	fmt.Println(res.String())
	return nil
}

func (s *Storage) Init(ctx context.Context) error {

	return nil
}
