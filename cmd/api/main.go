package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"restaurant-backend/internal/auth"
	"restaurant-backend/internal/database"
)

func main() {
	// ============================================================
	// 1. Load .env
	// ============================================================
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using env vars")
	}

	// ============================================================
	// 2. Database
	// ============================================================
	db, err := database.Connect() // also loads .env
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database connected successfully!")

	// ============================================================
	// 3. Auth Module
	// ============================================================
	authCfg, err := auth.LoadConfig()
	if err != nil {
		log.Fatalf("Auth config error: %v", err)
	}
	authModule := auth.NewModule(auth.ModuleDeps{
		DB:       db,
		Branches: branchChecker{db: db},
		Config:   authCfg,
	})

	seedCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	err = authModule.Service.SeedCEO(seedCtx,
		os.Getenv("SEED_CEO_NAME"), os.Getenv("SEED_CEO_EMAIL"), os.Getenv("SEED_CEO_PASSWORD"))
	cancel()
	if err != nil {
		log.Fatalf("CEO seed failed: %v", err)
	}
	log.Println("✅ Auth module ready")

	// ============================================================
	// 4. Gin Router
	// ============================================================
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	// Health (public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "auth-api",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1
	api := router.Group("/api/v1")
	authModule.Handler.RegisterRoutes(api)

	log.Println("✅ Auth routes registered")

	// ============================================================
	// 5. Start Server
	// ============================================================
	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}

	log.Printf("🚀 Auth API listening on :%s", port)
	log.Printf("   ❤️  Health:   http://localhost:%s/health", port)
	log.Printf("   🔗 BaseURL:  http://localhost:%s/api/v1", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
