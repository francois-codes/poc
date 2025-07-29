package main

import (
	"context"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func main() {

	db := InitDB()
	defer db.Close()

	app := fiber.New()

	// Route simple avec accès à la DB
	app.Get("/ping", func(c *fiber.Ctx) error {
		var now string
		err := db.QueryRow(c.Context(), "SELECT NOW()").Scan(&now)
		if err != nil {
			return c.Status(500).SendString("Erreur DB: " + err.Error())
		}
		return c.JSON(fiber.Map{"time": now})
	})

	log.Println("🚀 Serveur démarré sur http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
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
