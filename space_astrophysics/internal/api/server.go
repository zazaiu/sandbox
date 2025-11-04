package api

import (
	"log"
	"net/http"
	_ "space_astrophysics/internal/api/docs"
	"space_astrophysics/internal/app/handler"
	"space_astrophysics/internal/app/models"
	"space_astrophysics/internal/app/repository"
	"space_astrophysics/internal/utils"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	log.Println("✅ Database initialized")
}

// ----------------------
// ИНИЦИАЛИЗАЦИЯ ТЕСТОВЫХ ПОЛЬЗОВАТЕЛЕЙ
// ----------------------
func InitTestUsers(db *gorm.DB) {
	repo := repository.NewRepository(db)

	// Создаем обычного пользователя
	user := &models.User{
		Username: "user",
		Password: "1234",
		Role:     "astronaut",
	}

	// Создаем модератора
	moderator := &models.User{
		Username: "moderator",
		Password: "1234",
		Role:     "mission_control",
	}

	// Создаем пользователей (игнорируем ошибки если уже существуют)
	if err := repo.CreateUser(user); err != nil {
		log.Printf("⚠️ User creation: %v", err)
	} else {
		log.Printf("✅ Created user: %s", user.Username)
	}

	if err := repo.CreateUser(moderator); err != nil {
		log.Printf("⚠️ Moderator creation: %v", err)
	} else {
		log.Printf("✅ Created moderator: %s", moderator.Username)
	}

	log.Println("🎉 Test users initialized: user/1234 (astronaut), moderator/1234 (mission_control)")
}

// ----------------------
// ЗАПУСК СЕРВЕРА
// ----------------------
func StartServer() {
	InitDB()

	// Инициализируем тестовых пользователей
	InitTestUsers(DB)

	repo := repository.NewRepository(DB)
	h := handler.NewHandler(repo)

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	utils.InitRedis("localhost:6379", "", 0)

	// -------------------------------
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
	// Swagger UI
	// -------------------------------
	//docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// -------------------------------
	// REST API
	// -------------------------------

	// === ГОСТЬ (без токена) ===
	api := r.Group("/api")
	{
		api.POST("/users/register", h.RegisterUser)
		api.POST("/users/login", h.Login)
		api.GET("/planets", h.ListPlanets)
		api.GET("/planets/:id", h.ShowPlanetDetail)
	}

	// === Авторизованные пользователи (любая роль) ===
	auth := api.Group("/")
	auth.Use(h.AuthMiddleware()) // Без указания ролей - проверяет только аутентификацию
	{
		// Личный кабинет
		auth.GET("/users/profile", h.GetUserProfile)
		auth.PUT("/users/profile", h.UpdateUserProfile)
		auth.POST("/users/logout", h.Logout)

		// Заявки
		auth.GET("/worlds", h.ListWorldsFiltered)
		auth.GET("/worlds/:id", h.ViewWorld)
		auth.PUT("/worlds/:id", h.UpdateWorld)
		auth.DELETE("/worlds/:id", h.DeleteWorld)

		// Корзина
		auth.GET("/cart", h.GetCartIcon)
		auth.POST("/worlds/:id/planet/:planet_id", h.AddPlanetToWorld)
	}

	// === Только для астронавтов ===
	astronaut := api.Group("/")
	astronaut.Use(h.AuthMiddleware("astronaut"))
	{
		// Оформление заявки
		astronaut.PUT("/worlds/:id/form", h.FormWorld)
	}

	// === Только для модератора (mission_control) ===
	admin := api.Group("/")
	admin.Use(h.AuthMiddleware("mission_control"))
	{
		// Управление планетами
		admin.POST("/planets", h.CreatePlanet)
		admin.PUT("/planets/:id", h.UpdatePlanet)
		admin.DELETE("/planets/:id", h.DeletePlanet)
		admin.POST("/planets/:id/image", h.UploadPlanetImage)

		// Завершение заявок
		admin.PUT("/worlds/:id/complete", h.CompleteWorld)

		// Работа с связями M-M - ПЕРЕНОСИМ В ОТДЕЛЬНУЮ ГРУППУ
		admin.GET("/world-planets/:world_id/:planet_id", h.GetWorldPlanet)
		admin.DELETE("/world-planets/:world_id/:planet_id", h.DeleteWorldPlanet)
		admin.PUT("/world-planets/:world_id/:planet_id", h.UpdateWorldPlanet)

		// Админка Redis
		admin.GET("/admin/sessions", h.ShowAllSessions)
		admin.GET("/admin/redis-check", h.RedisCheck)
	}

	// -------------------------------
	// ЗАПУСК СЕРВЕРА
	// -------------------------------
	log.Println("🚀 Server running on http://localhost:8080")
	log.Println("📚 Swagger available at http://localhost:8080/swagger/index.html")
	log.Println("🔐 Admin sessions at http://localhost:8080/api/admin/sessions")
	log.Println("👤 Test users: user/1234 (astronaut), moderator/1234 (mission_control)")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
