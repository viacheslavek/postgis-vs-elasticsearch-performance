package app

import (
	"context"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/benchmark"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/elasticsearch"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage/postgis"
	"log"
)

// TODO: оформить ридми с описанием запуска
// TODO: ридми с замерами (когда сделаю возможность запускать N тестов)

// TODO: пора бы и тесты написать - после основных дополнений

func Run(pg *postgis.Storage, es *elasticsearch.Storage) {

	log.Printf("app start\n")

	ctx := context.Background()

	if err := benchmark.RunBenchPoint(ctx, pg, "postgis", 2e6, 20, 10); err != nil {
		log.Fatalf("can't do pg bench %e\n", err)
	}

	if err := benchmark.RunBenchPoint(ctx, es, "elasticseacrh", 2e6, 20, 10); err != nil {
		log.Fatalf("can't do es bench %e\n", err)
	}

	if err := benchmark.RunBenchPolygon(ctx, pg, "postgis", 2e6); err != nil {
		log.Fatalf("can't do pg polygon bench %e\n", err)
	}

	if err := benchmark.RunBenchPolygon(ctx, es, "elasticseacrh", 2e6); err != nil {
		log.Fatalf("can't do es polygon bench %e\n", err)
	}

	log.Printf("end\n")

}
