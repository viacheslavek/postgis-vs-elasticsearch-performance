package app

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/benchmark"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

func Run(pg *postgis.Storage, es *elasticsearch.Storage) {

	log.Printf("app start\n")

	ctx := context.Background()

	bfPg, err := benchmark.RunBenchNCheck(ctx, pg, pg, "postgis", 2e2, 20, 10, 25)

	if err != nil {
		log.Fatalf("can't do pg bench %e\n", err)
	}

	fmt.Println(bfPg)

	bfEs, err := benchmark.RunBenchNCheck(ctx, es, es, "elasticsearch", 2e2, 20, 10, 25)

	if err != nil {
		log.Fatalf("can't do es bench %e\n", err)
	}

	fmt.Println(bfEs)

	log.Printf("end\n")

}
