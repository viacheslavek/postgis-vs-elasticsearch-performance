package postgis

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal/storage"
	"github.com/jackc/pgx/v5"
	"log"
)

type Storage struct {
	db *pgx.Conn
}

func New() (*Storage, error) {
	//databaseURL := os.Getenv("DATABASE_URL")
	//
	//if databaseURL == "" {
	//	return nil, fmt.Errorf("space db url\n")
	//}
	//
	//db, err := pgx.Connect(context.Background(), databaseURL)

	db, err := pgx.Connect(context.Background(), "postgres://slava:passwordforgis@localhost:5432/postgresgis")

	if err != nil {
		return nil, fmt.Errorf("unable to connect database %w\n", err)
	}

	var greeting string
	err = db.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return nil, fmt.Errorf("QueryRow failed: %w\n", err)
	}

	log.Printf("connect db %s\n", greeting)

	return &Storage{db}, nil
}

// Перенести это в миграции
// Можно проверить работоспособность с индексами по точкам и без
// CREATE INDEX moscow_region_geom_idx ON moscow_region USING GIST (geom);

func (s *Storage) Init(ctx context.Context) error {
	q := `
		CREATE EXTENSION IF NOT EXISTS postgis;
		CREATE TABLE IF NOT EXISTS moscow_region (
    		id SERIAL PRIMARY KEY ,
    		geom GEOMETRY(Point, 4326)
		);
	`

	_, err := s.db.Exec(
		ctx,
		q,
	)
	if err != nil {
		return fmt.Errorf("can't create tables %w", err)
	}

	return nil
}

func (s *Storage) AddPoint(ctx context.Context, p storage.Point) error {

	q := `
		INSERT INTO moscow_region (geom) 
			VALUES (ST_SetSRID(ST_MakePoint($1, $2), 4326))
	`

	_, err := s.db.Query(ctx, q, p.Longitude, p.Latitude)

	if err != nil {
		return fmt.Errorf("can't add a point %e\n", err)
	}

	return nil
}
