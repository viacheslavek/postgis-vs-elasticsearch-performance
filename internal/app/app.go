package app

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/benchmark"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

func Run(pg *postgis.Storage, es *elasticsearch.Storage) {

	log.Printf("app start\n")

	ctx := context.Background()

	if err := benchmark.RunBenchPoint(ctx, pg, "postgis", 20e1, 5, 10); err != nil {
		log.Fatalf("can't do pg bench %e\n", err)
	}

	if err := benchmark.RunBenchPoint(ctx, es, "elasticseacrh", 20e1, 5, 10); err != nil {
		log.Fatalf("can't do es bench %e\n", err)
	}

	if err := benchmark.RunBenchPolygon(ctx, pg, "postgis", 20e1); err != nil {
		log.Fatalf("can't do pg polygon bench %e\n", err)
	}

	if err := benchmark.RunBenchPolygon(ctx, es, "elasticseacrh", 20e1); err != nil {
		log.Fatalf("can't do es polygon bench %e\n", err)
	}

	log.Printf("end\n")

}
