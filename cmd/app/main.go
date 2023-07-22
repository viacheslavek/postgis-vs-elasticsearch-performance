package main

import (
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/app"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

func main() {

	log.Printf("app starts launching\n")

	es, err := elasticsearch.New()
	if err != nil {
		log.Fatalf("can't connect to elasticsearch %e\n", err)
	}

	pg, err := postgis.New()
	if err != nil {
		log.Fatalf("can't connect to elasticsearch %e\n", err)
	}

	app.Run(pg, es)
}
