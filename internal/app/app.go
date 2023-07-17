package app

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
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

	p := internal.Point{Latitude: 37.190000, Longitude: 51.100000}

	points, err := es.GetInRadius(ctx, p, 10000000)

	if err != nil {
		fmt.Println("опять 25", err)
	}

	fmt.Println(points)

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
