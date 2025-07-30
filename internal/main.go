package main

import (
	"cognyx/psychic-robot/json"
	"cognyx/psychic-robot/persistence/db"
	"cognyx/psychic-robot/persistence/repository"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strconv"
	"time"
)

const GENERATE = false

func main() {

	dbconn := InitDB()
	defer dbconn.Close()

	if GENERATE {
		db.Generate(dbconn)
	}
	// Repository
	queries := db.New(dbconn) // 👈 conversion pool → Queries
	repo := repository.NewVersionRepository(queries)

	app := fiber.New()
	app.Get("/datamodel/:id", func(c *fiber.Ctx) error {
		t := time.Now()
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalide"})
		}

		jsonData, err := repo.GetLatestByDatamodelID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erreur base de données"})
		}
		log.Println(time.Since(t))
		return c.Type("application/json").Send(jsonData)
	})

	app.Get("/datamodel/update/:id", func(c *fiber.Ctx) error {
		t := time.Now()
		idStr := c.Params("id")
		_, err := strconv.ParseInt(idStr, 10, 64) //TODO: change for DB JSON
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id invalide"})
		}

		jsonData, _ := json.CompareJSONFiles("/Users/thomas/go/poc/json/tree.json", "/Users/thomas/go/poc/json/tree2.json")

		log.Println(time.Since(t))
		return c.Type("application/json").Send([]byte(jsonData))
	})

	log.Println("🚀 Server started on http://localhost:8585")
	log.Fatal(app.Listen(":8585"))
}

func InitDB() *pgxpool.Pool {
	dsn := "postgres://cognyx:cognyx@localhost:5432/cognyx"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Erreur de création du pool PostgreSQL : %v", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("PostgreSQL inaccessible : %v", err)
	}

	return pool
}
