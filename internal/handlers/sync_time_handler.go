package handlers

import (
	"os/exec"
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type SyncTimeHandler struct {
	cfg *config.Config
}

func NewSyncTimeHandler(cfg *config.Config) *SyncTimeHandler {
	return &SyncTimeHandler{cfg: cfg}
}

type SyncTimeRequest struct {
	Timestamp string `json:"timestamp"`
}

func (h *SyncTimeHandler) SyncTime(c *fiber.Ctx) error {
	var request SyncTimeRequest
	if err := c.BodyParser(&request); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if request.Timestamp == "" {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Time cannot be empty")
	}

	t, err := time.Parse(time.RFC3339, request.Timestamp)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid timestamp format")
	}

	formattedTime := t.Format("2006-01-02 15:04:05")
	cmd := exec.Command("date", "-s", formattedTime)
	err = cmd.Run()
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Cannot sync time")
	}

	log.Info("Time synced to", formattedTime)

	return helpers.SuccessResponse(c,
		fiber.StatusOK,
		"Time synced",
		nil,
	)
}
