package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tidwall/randjson"
	"log"
	"sync"
	"time"
)

const (
	totalDatamodels      = 1000
	versionsPerDatamodel = 10000
	numWorkers           = 100
)

func Generate(dbpool *pgxpool.Pool) {
	ctx := context.Background()
	log.Printf("🚀 Inserting %d datamodels + %d versions using %d workers...\n", totalDatamodels, totalDatamodels*versionsPerDatamodel, numWorkers)

	start := time.Now()
	datamodelCh := make(chan int, totalDatamodels)

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			q := New(dbpool)

			for i := range datamodelCh {
				insertDatamodelWithVersions(ctx, dbpool, q, i)
			}
		}(w)
	}

	// Feed datamodel indexes to workers
	for i := 0; i < totalDatamodels; i++ {
		datamodelCh <- i
	}
	close(datamodelCh)

	wg.Wait()

	log.Printf("✅ Done in %s", time.Since(start))
}

func insertDatamodelWithVersions(ctx context.Context, dbpool *pgxpool.Pool, q *Queries, index int) error {
	tx, err := dbpool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := q.WithTx(tx)

	// Step 1: insert datamodel
	dm, _ := qtx.CreateDatamodel(ctx, fmt.Sprintf("Model #%d", index+1))

	for v := 1; v <= versionsPerDatamodel; v++ {
		randomJSON := randjson.Make(12, nil)

		version := CreateVersionParams{
			ObjectType: "datamodel",
			ObjectID:   dm.ID,
			Json:       randomJSON,
			Version:    int32(v),
			Action:     "update",
			Actor:      "system",
		}

		qtx.CreateVersion(ctx, version)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	log.Printf("✅ Completed datamodel %d", index+1)
	return nil
}

func Generate2(dbpool *pgxpool.Pool) {
	ctx := context.Background()
	q := New(dbpool)

	const totalDatamodels = 1000
	const versionsPerDatamodel = 10000
	const totalVersions = totalDatamodels * versionsPerDatamodel

	log.Printf("🚀 Inserting %d datamodels + %d versions...\n", totalDatamodels, totalVersions)

	start := time.Now()

	tx, err := dbpool.Begin(ctx)
	if err != nil {
		log.Fatal("begin tx:", err)
	}
	defer tx.Rollback(ctx) // in case of panic

	qtx := q.WithTx(tx)

	for i := 0; i < totalDatamodels; i++ {
		// Step 1: insert datamodel
		dm, err := qtx.CreateDatamodel(ctx, fmt.Sprintf("Model #%d", i+1))
		if err != nil {
			log.Fatalf("insert datamodel %d: %v", i+1, err)
		}

		for v := 1; v <= versionsPerDatamodel; v++ {
			randomJSON := randjson.Make(12, nil)

			version := CreateVersionParams{
				ObjectType: "datamodel",
				ObjectID:   dm.ID,
				Json:       randomJSON,
				Version:    int32(v),
				Action:     "update",
				Actor:      "system",
			}

			if _, err := qtx.CreateVersion(ctx, version); err != nil {
				log.Fatalf("insert version (%d-%d): %v", i+1, v, err)
			}

			if v%1000 == 0 {
				log.Printf("Inserted %d versions for datamodel %d", v, i+1)
			}
		}

		if (i+1)%10 == 0 {
			log.Printf("✅ %d datamodels processed", i+1)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal("commit tx:", err)
	}

	log.Printf("✅ Done in %s", time.Since(start))
}
