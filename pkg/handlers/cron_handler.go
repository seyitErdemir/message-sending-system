package handlers

import (
	"fiber-app/pkg/cron"
	"fiber-app/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type CronStatusResponse struct {
	Status    string `json:"status" example:"success"`
	IsRunning bool   `json:"is_running" example:"true"`
}

type CronMessageResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Cron started"`
}

type CronLogsResponse struct {
	Status string           `json:"status" example:"success"`
	Data   []models.CronLog `json:"data"`
}

// @Summary Start cron job
// @Description Starts the message sending cron job
// @Tags cron
// @Accept json
// @Produce json
// @Success 200 {object} CronMessageResponse "Successful response"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /cron/start [post]
func StartCronJob(c *fiber.Ctx) error {
	if err := cron.StartCron(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Failed to start cron job",
			Code:    "CRON_START_ERROR",
		})
	}

	return c.JSON(CronMessageResponse{
		Status:  "success",
		Message: "Cron job started",
	})
}

// @Summary Stop cron job
// @Description Stops the message sending cron job
// @Tags cron
// @Accept json
// @Produce json
// @Success 200 {object} CronMessageResponse "Successful response"
// @Router /cron/stop [post]
func StopCronJob(c *fiber.Ctx) error {
	cron.StopCron()
	return c.JSON(CronMessageResponse{
		Status:  "success",
		Message: "Cron job stopped",
	})
}

// @Summary Get cron job status
// @Description Checks if the cron job is running
// @Tags cron
// @Accept json
// @Produce json
// @Success 200 {object} CronStatusResponse "Successful response"
// @Router /cron/status [get]
func GetCronStatus(c *fiber.Ctx) error {
	return c.JSON(CronStatusResponse{
		Status:    "success",
		IsRunning: cron.IsCronRunning(),
	})
}

// @Summary Get cron logs
// @Description Retrieves the cron job execution logs
// @Tags cron
// @Accept json
// @Produce json
// @Success 200 {object} CronLogsResponse "Successful response"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /cron/logs [get]
func GetCronLogs(c *fiber.Ctx) error {
	logs, err := cron.GetCronLogs(100)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Failed to retrieve cron logs",
			Code:    "CRON_LOGS_ERROR",
		})
	}

	return c.JSON(CronLogsResponse{
		Status: "success",
		Data:   logs,
	})
}
