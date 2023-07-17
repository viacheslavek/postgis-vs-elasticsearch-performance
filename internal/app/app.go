package app

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"log"
)

func Run(es *elasticsearch.Storage) {

	log.Printf("app start\n")

	ctx := context.Background()

	err := es.Init(ctx)
	if err != nil {
		log.Fatalf("can't init es %e\n", err)
	}

	log.Printf("connect to es\n")

	fmt.Println("end")

}

// curl -X GET "http://localhost:9200/moscow_region/_search" -H "Content-Type: application/json" -d'
//{
//  "query": {
//    "match_all": {}
//  },
//  "size": 100
//}
//'
