package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"space_astrophysics/internal/app/models"
	"space_astrophysics/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repo: r}
}

const currentUserID = 1 // имитация авторизации (создатель)

// =========================================================
// ВСПОМОГАТЕЛЬНЫЕ
// =========================================================
func (h *Handler) getOrCreateDraftWorld(userID int) (*models.World, error) {
	draft, err := h.Repo.GetDraftWorld(userID)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		newWorld := &models.World{
			CreatorID:   userID,
			WorldStatus: "draft",
			CreatedAt:   time.Now(),
		}
		if err := h.Repo.CreateWorld(newWorld); err != nil {
			return nil, err
		}
		draft = newWorld
	}
	return draft, nil
}

// =========================================================
// 🪐 PLANETS (Услуги)
// =========================================================

// GET /api/planets?q=filter
func (h *Handler) ListPlanets(ctx *gin.Context) {
	q := ctx.Query("q")
	planets, err := h.Repo.GetAllPlanets()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var filtered []models.Planet
	for _, p := range planets {
		if q == "" || strings.Contains(strings.ToLower(p.Name), strings.ToLower(q)) {
			filtered = append(filtered, p)
		}
	}
	ctx.JSON(http.StatusOK, filtered)
}

// GET /api/planets/:id
func (h *Handler) ShowPlanetDetail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	planet, err := h.Repo.GetPlanetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Планета не найдена"})
		return
	}
	ctx.JSON(http.StatusOK, planet)
}

// POST /api/planets
func (h *Handler) CreatePlanet(ctx *gin.Context) {
	var p models.Planet
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}
	p.Status = "active"
	if err := h.Repo.CreatePlanet(&p); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, p)
}

// PUT /api/planets/:id
func (h *Handler) UpdatePlanet(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var update models.Planet
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.UpdatePlanet(id, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Планета обновлена"})
}

// DELETE /api/planets/:id
func (h *Handler) DeletePlanet(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.Repo.DeletePlanet(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Планета удалена"})
}

// POST /api/planets/:id/image
func (h *Handler) UploadPlanetImage(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Файл не найден"})
		return
	}
	filename, err := h.Repo.UploadPlanetImageToMinio(id, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Изображение обновлено", "url": filename})
}

// POST /api/planets/:id/add-to-world
func (h *Handler) AddPlanetToWorld(ctx *gin.Context) {
	planetID, _ := strconv.Atoi(ctx.Param("id"))
	world, err := h.getOrCreateDraftWorld(currentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.AddPlanetToWorld(world.ID, planetID, 1, false); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Планета добавлена в черновик"})
}

// =========================================================
// 🌍 WORLDS (Заявки)
// =========================================================

// GET /api/cart
func (h *Handler) GetCartIcon(ctx *gin.Context) {
	world, err := h.Repo.GetDraftWorld(currentUserID)
	if err != nil || world == nil {
		ctx.JSON(http.StatusOK, gin.H{"world_id": nil, "count": 0})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"world_id": world.ID, "count": len(world.Planets)})
}

// GET /api/worlds?status=&from=&to=
func (h *Handler) ListWorldsFiltered(ctx *gin.Context) {
	status := ctx.Query("status")
	from := ctx.Query("from")
	to := ctx.Query("to")
	worlds, err := h.Repo.GetWorldsFiltered(status, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, worlds)
}

// GET /api/worlds/:id
func (h *Handler) ViewWorld(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	world, err := h.Repo.GetWorldByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}
	ctx.JSON(http.StatusOK, world)
}

// PUT /api/worlds/:id
func (h *Handler) UpdateWorld(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var update map[string]interface{}
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	delete(update, "id")
	delete(update, "creator_id")
	delete(update, "world_status")
	if err := h.Repo.UpdateWorldFields(id, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка обновлена"})
}

// PUT /api/worlds/:id/form
func (h *Handler) FormWorld(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.Repo.FormWorld(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка оформлена"})
}

// PUT /api/worlds/:id/complete
func (h *Handler) CompleteWorld(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.Repo.CompleteWorld(id, 2); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка завершена"})
}

// DELETE /api/worlds/:id
func (h *Handler) DeleteWorld(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.Repo.DeleteWorldSQL(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка удалена"})
}

// =========================================================
// ⚙️ M-M связи (WorldPlanet)
// =========================================================

// DELETE /api/worlds/:world_id/planet/:planet_id
func (h *Handler) DeleteWorldPlanet(ctx *gin.Context) {
	worldID, _ := strconv.Atoi(ctx.Param("world_id"))
	planetID, _ := strconv.Atoi(ctx.Param("planet_id"))
	if err := h.Repo.DeleteWorldPlanet(worldID, planetID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Планета удалена из заявки"})
}

// PUT /api/worlds/:world_id/planet/:planet_id
func (h *Handler) UpdateWorldPlanet(ctx *gin.Context) {
	worldID, _ := strconv.Atoi(ctx.Param("world_id"))
	planetID, _ := strconv.Atoi(ctx.Param("planet_id"))
	var payload struct {
		Quantity int  `json:"quantity"`
		IsMain   bool `json:"is_main"`
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}
	if err := h.Repo.UpdateWorldPlanet(worldID, planetID, payload.Quantity, payload.IsMain); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Связь обновлена"})
}

// =========================================================
// 👤 USERS
// =========================================================

// POST /api/users/register
func (h *Handler) RegisterUser(ctx *gin.Context) {
	var u models.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.CreateUser(&u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, u)
}

// GET /api/users/me
func (h *Handler) GetUserProfile(ctx *gin.Context) {
	user, err := h.Repo.GetUserByID(currentUserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// PUT /api/users/me
func (h *Handler) UpdateUserProfile(ctx *gin.Context) {
	var update models.User
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.UpdateUser(currentUserID, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Профиль обновлен"})
}

// POST /api/users/login
func (h *Handler) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Repo.Authenticate(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Успешный вход", "user": user})
}

// POST /api/users/logout
func (h *Handler) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь вышел"})
}

// =========================================================
// Утилиты: расчёт орбиты
// =========================================================
func julianDate(t time.Time) float64 {
	year, month, day := t.Date()
	if month <= 2 {
		year--
		month += 12
	}
	A := year / 100
	B := 2 - A + A/4
	return float64(int(365.25*float64(year+4716))) +
		float64(int(30.6001*float64(month+1))) +
		float64(day) + float64(B) - 1524.5
}

func solveKepler(M, e float64) float64 {
	E := M
	for i := 0; i < 15; i++ {
		E = E - (E-e*math.Sin(E)-M)/(1-e*math.Cos(E))
	}
	return E
}

func calcOrbit(p models.Planet, t time.Time) (float64, float64) {
	jd := julianDate(t)
	M := 2 * math.Pi * (jd - p.T0) / (p.Period * 365.25)
	M = math.Mod(M, 2*math.Pi)
	E := solveKepler(M, p.E)
	nu := 2 * math.Atan2(
		math.Sqrt(1+p.E)*math.Sin(E/2),
		math.Sqrt(1-p.E)*math.Cos(E/2),
	)
	r := p.A * (1 - p.E*math.Cos(E))
	return r, nu * 180 / math.Pi
}
