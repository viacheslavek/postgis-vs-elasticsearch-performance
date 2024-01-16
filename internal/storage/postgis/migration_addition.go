package postgis

import (
	"context"
	"fmt"
	"github.com/VyacheslavIsWorkingNow/postgis-vs-elasticsearch-performance/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Storage struct {
	db   *pgx.Conn
	pool *pgxpool.Pool
}

func New() (*Storage, error) {

	db, err := pgx.Connect(context.Background(), "postgres://slava:passwordforgis@localhost:5432/postgresgis")

	if err != nil {
		return nil, fmt.Errorf("unable to connect database %w\n", err)
	}

	var greeting string
	err = db.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return nil, fmt.Errorf("QueryRow failed: %w\n", err)
	}

	config, err := pgxpool.ParseConfig("postgres://slava:passwordforgis@localhost:5432/postgresgis")
	if err != nil {
		return nil, fmt.Errorf("can`t parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("can`t create connection pool: %w", err)
	}

	log.Printf("connect db %s\n", greeting)

	return &Storage{db, pool}, nil
}

func (s *Storage) Close(ctx context.Context) error {

	err := s.db.Close(ctx)
	if err != nil {
		return fmt.Errorf("can't close connection %e\n", err)
	}
	s.pool.Close()

	return nil
}

// CREATE INDEX moscow_region_geom_idx ON moscow_region USING GIST (geom);

func (s *Storage) Init(ctx context.Context) error {
	q := `
		CREATE EXTENSION IF NOT EXISTS postgis;
		CREATE TABLE IF NOT EXISTS moscow_region (
    		id SERIAL PRIMARY KEY,
    		geom GEOMETRY(Point, 4326)
		);
		CREATE INDEX moscow_region_geom_idx ON moscow_region USING GIST (geom);
	`

	return s.initBase(ctx, q)
}

func (s *Storage) Drop(ctx context.Context) error {
	q := `
		DROP TABLE IF EXISTS moscow_region;
	`

	return s.drop(ctx, q)
}

func (s *Storage) AddPoint(ctx context.Context, p internal.Point) error {

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

func (s *Storage) AddPointBatch(ctx context.Context, points []internal.Point) error {

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("can`t acquire connection from pool: %w\n", err)
	}

	defer conn.Release()

	batch := &pgx.Batch{}

	q := `
		INSERT INTO moscow_region (geom) 
			VALUES (ST_SetSRID(ST_MakePoint($1, $2), 4326))
	`

	for _, p := range points {
		batch.Queue(q, p.Longitude, p.Latitude)
	}

	return s.addBatch(ctx, conn, batch, len(points))
}

func (s *Storage) DeletePoint(ctx context.Context, p internal.Point) error {

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

func (s *Storage) IsPointExist(ctx context.Context, p internal.Point) (bool, error) {

	q := `
		SELECT count(*) AS count
		FROM moscow_region
		WHERE geom = ST_SetSRID(ST_MakePoint($1, $2), 4326)
	`

	row, err := s.db.Query(ctx, q, p.Longitude, p.Latitude)

	if err != nil {
		return false, fmt.Errorf("can't delete a point %e\n", err)
	}

	defer row.Close()

	if !row.Next() {
		return true, nil
	}

	return false, nil
}

func (s *Storage) initBase(ctx context.Context, query string) error {
	_, err := s.db.Exec(
		ctx,
		query,
	)
	if err != nil {
		return fmt.Errorf("can't create tables %w", err)
	}

	return nil
}

func (s *Storage) drop(ctx context.Context, query string) error {
	_, err := s.db.Exec(
		ctx,
		query,
	)
	if err != nil {
		return fmt.Errorf("can't drop tables %w", err)
	}

	return nil
}

func (s *Storage) addBatch(ctx context.Context, conn *pgxpool.Conn, batch *pgx.Batch, N int) error {
	results := conn.SendBatch(ctx, batch)

	// необязательно, если мы хотим максимальную скорость
	for i := 0; i < N; i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("can't execute batch queure: %w\n", err)
		}
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("can`t close batch results: %w", err)
	}

	return nil
}
