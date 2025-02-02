package handlers

import (
	"fiber-app/pkg/cron"
	"fiber-app/pkg/models"

	"github.com/gofiber/fiber/v2"
)

// CronStatusResponse represents the cron status response
type CronStatusResponse struct {
	Status    string `json:"status" example:"success"`  // İşlem durumu
	IsRunning bool   `json:"is_running" example:"true"` // Cron job çalışma durumu
}

// CronMessageResponse represents the cron operation message response
type CronMessageResponse struct {
	Status  string `json:"status" example:"success"`          // İşlem durumu
	Message string `json:"message" example:"Cron başlatıldı"` // İşlem mesajı
}

// CronLogsResponse represents the cron logs response
type CronLogsResponse struct {
	Status string           `json:"status" example:"success"` // İşlem durumu
	Data   []models.CronLog `json:"data"`                     // Cron logları
}

// @Summary Cron job'ı başlat
// @Description Mesaj gönderme cron job'ını başlatır
// @Tags cron
// @Accept json
// @Produce json
// @Success 200 {object} CronMessageResponse "Başarılı yanıt"
// @Failure 500 {object} ErrorResponse "Sunucu hatası"
// @Router /cron/start [post]
func StartCronJob(c *fiber.Ctx) error {
	if err := cron.StartCron(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Cron job başlatılamadı",
			Code:    "CRON_START_ERROR",
		})
	}

	return c.JSON(CronMessageResponse{
		Status:  "success",
		Message: "Cron job başlatıldı",
	})
}

// @Summary Cron job'ı durdur
// @Description Mesaj gönderme cron job'ını durdurur
// @Tags cron
// @Accept json
// @Produce json
// @Success 200 {object} CronMessageResponse "Başarılı yanıt"
// @Router /cron/stop [post]
func StopCronJob(c *fiber.Ctx) error {
	cron.StopCron()
	return c.JSON(CronMessageResponse{
		Status:  "success",
		Message: "Cron job durduruldu",
	})
}

// @Summary Cron job durumunu getir
// @Description Cron job'ın çalışıp çalışmadığını kontrol eder
// @Tags cron
// @Accept json
// @Produce json
// @Success 200 {object} CronStatusResponse "Başarılı yanıt"
// @Router /cron/status [get]
func GetCronStatus(c *fiber.Ctx) error {
	return c.JSON(CronStatusResponse{
		Status:    "success",
		IsRunning: cron.IsCronRunning(),
	})
}

// @Summary Cron loglarını getir
// @Description Cron job'ın çalışma loglarını getirir
// @Tags cron
// @Accept json
// @Produce json
// @Success 200 {object} CronLogsResponse "Başarılı yanıt"
// @Failure 500 {object} ErrorResponse "Sunucu hatası"
// @Router /cron/logs [get]
func GetCronLogs(c *fiber.Ctx) error {
	logs, err := cron.GetCronLogs(100) // Son 100 log
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Cron logları getirilemedi",
			Code:    "CRON_LOGS_ERROR",
		})
	}

	return c.JSON(CronLogsResponse{
		Status: "success",
		Data:   logs,
	})
}
