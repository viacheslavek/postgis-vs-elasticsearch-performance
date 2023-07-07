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

func (s *Storage) Close(ctx context.Context) error {

	err := s.db.Close(ctx)
	if err != nil {
		return fmt.Errorf("can't close connection %e\n", err)
	}

	return nil
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

func (s *Storage) Drop(ctx context.Context) error {
	q := `
		DROP TABLE IF EXISTS moscow_region;
	`

	_, err := s.db.Exec(
		ctx,
		q,
	)
	if err != nil {
		return fmt.Errorf("can't drop tables %w", err)
	}

	return nil
}

func (s *Storage) AddPoint(ctx context.Context, p storage.Point) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return fmt.Errorf("can't begin transaction: %w\n", err)
	}

	q := `
		INSERT INTO moscow_region (geom) 
			VALUES (ST_SetSRID(ST_MakePoint($1, $2), 4326))
	`

	_, err = tx.Exec(ctx, q, p.Longitude, p.Latitude)

	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("can't add a point: %w\n", err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return fmt.Errorf("can't commit transactions %w\n", err)
	}

	return nil
}

func (s *Storage) DeletePoint(ctx context.Context, p storage.Point) error {

	q := `
		DELETE FROM moscow_region
		WHERE geom = ST_SetSRID(ST_MakePoint($1, $2), 4326)
	`

	_, err := s.db.Query(ctx, q, p.Longitude, p.Latitude)

	if err != nil {
		return fmt.Errorf("can't delete a point %e\n", err)
	}

	return nil
}

func (s *Storage) IsPoint(ctx context.Context, p storage.Point) (bool, error) {

	q := `
		SELECT count(*) AS count
		FROM moscow_region
		WHERE geom = ST_SetSRID(ST_MakePoint($1, $2), 4326)
	`

	row, err := s.db.Query(ctx, q, p.Longitude, p.Latitude)

	defer row.Close()

	if err != nil {
		return false, fmt.Errorf("can't delete a point %e\n", err)
	}

	if row.Next() != false {
		return true, nil
	}

	return false, nil
}
