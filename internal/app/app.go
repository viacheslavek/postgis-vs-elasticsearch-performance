package app

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

func Run(pg *postgis.Storage, es *elasticsearch.Storage) {

	log.Printf("app start\n")
	ctx := context.Background()

	if err := es.Init(ctx); err != nil {
		log.Fatalf("can't init es %e\n", err)
	}

	if err := pg.Init(ctx); err != nil {
		log.Fatalf("can't init pg %e\n", err)
	}

	log.Printf("start benchmark\n")

	// Запуск различных бенчев

	log.Printf("end\n")

}

// curl -X GET "http://localhost:9200/moscow_region/_search" -H "Content-Type: application/json" -d'
//{
//  "query": {
//    "match_all": {}
//  },
//  "size": 100
//}
//'
