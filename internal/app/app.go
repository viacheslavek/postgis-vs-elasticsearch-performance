package app

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/benchmark"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

const (
	benchPathPG = "cmd/bench_result_postgis.html"
	benchPathES = "cmd/bench_result_elasticsearch.html"
)

func Run(pg *postgis.Storage, es *elasticsearch.Storage) {

	log.Printf("app start\n")

	ctx := context.Background()

	bfPg, err := benchmark.RunBenchNCheck(ctx, pg, pg, "postgis", 2e2, 200, 10, 15, 3)

	if err != nil {
		log.Fatalf("can't do pg bench %e\n", err)
	}

	fmt.Println(bfPg)

	if err = bfPg.ConvertToHTML(benchPathPG); err != nil {
		log.Fatalf("can't convert pg to html, %e", err)
	}

	bfEs, err := benchmark.RunBenchNCheck(ctx, es, es, "elasticsearch", 2e2, 200, 10, 15, 10)

	if err != nil {
		log.Fatalf("can't do es bench %e\n", err)
	}

	fmt.Println(bfEs)

	if err = bfEs.ConvertToHTML(benchPathES); err != nil {
		log.Fatalf("can't convert es to html, %e", err)
	}

	log.Printf("end\n")

}
