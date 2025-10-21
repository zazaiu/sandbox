package api

import (
	"log"
	"net/http"
	"space_astrophysics/internal/app/handler"
	"space_astrophysics/internal/app/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ----------------------
// ИНИЦИАЛИЗАЦИЯ БД
// ----------------------
func InitDB() {
	dsn := "host=localhost port=5433 user=astrouser password=1234 dbname=astrodb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}

	DB = db

	// Если таблицы уже созданы вручную
	log.Println("ℹ️ AutoMigrate пропущен (таблицы уже созданы вручную)")

	log.Println("✅ Database initialized")
}

// ----------------------
// ЗАПУСК СЕРВЕРА
// ----------------------
func StartServer() {
	InitDB()

	repo := repository.NewRepository(DB)
	h := handler.NewHandler(repo)

	r := gin.Default()

	// Статические файлы и шаблоны
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	// HTML фронт
	// -------------------------------
	r.GET("/planets", func(ctx *gin.Context) {
		planets, err := h.Repo.GetAllPlanets()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Ошибка: %v", err)
			return
		}
		ctx.HTML(http.StatusOK, "service_list.html", gin.H{
			"planets": planets,
		})
	})
	// -------------------------------
	// REST API
	// -------------------------------
	api := r.Group("/api")
	{
		// ---- PLANETS ----
		api.GET("/planets", h.ListPlanets)
		api.GET("/planets/:id", h.ShowPlanetDetail)
		api.POST("/planets", h.CreatePlanet)
		api.PUT("/planets/:id", h.UpdatePlanet)
		api.DELETE("/planets/:id", h.DeletePlanet)
		api.POST("/planets/:id/image", h.UploadPlanetImage)
		api.POST("/planets/:id/add-to-world", h.AddPlanetToWorld)

		// ---- WORLDS ----
		api.GET("/worlds", h.ListWorldsFiltered)
		api.GET("/worlds/:id", h.ViewWorld)
		api.PUT("/worlds/:id", h.UpdateWorld)
		api.PUT("/worlds/:id/form", h.FormWorld)
		api.PUT("/worlds/:id/complete", h.CompleteWorld)
		api.DELETE("/worlds/:id", h.DeleteWorld)

		api.GET("/cart", h.GetCartIcon)

		// ---- USERS ----
		api.POST("/users/register", h.RegisterUser)
		api.GET("/users/me", h.GetUserProfile)
		api.PUT("/users/me", h.UpdateUserProfile)
		api.POST("/users/login", h.Login)
		api.POST("/users/logout", h.Logout)

		// ---- WORLD-PLANET связи (M-M) ----
		wp := api.Group("/world-planets")
		{
			wp.DELETE("/:world_id/:planet_id", h.DeleteWorldPlanet)
			wp.PUT("/:world_id/:planet_id", h.UpdateWorldPlanet)
		}
		// --HTML--

	}

	log.Println("🚀 Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
